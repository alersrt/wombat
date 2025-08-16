package internal

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"log/slog"
	"strconv"
	"sync"
	"wombat/pkg/cipher"
)

var (
	botCommandRegister = tgbotapi.BotCommand{
		Command:     "register",
		Description: "RG",
	}
)

type TelegramSource struct {
	sourceType SourceType
	bot        *tgbotapi.BotAPI
	cipher     *cipher.AesGcmCipher
	router     *Router
	db         *DbStorage
	updChan    tgbotapi.UpdatesChannel
}

func NewTelegramSource(token string, router *Router, db *DbStorage, cipher *cipher.AesGcmCipher) (ts *TelegramSource, err error) {
	slog.Info("telegram:init:start")

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	updCfg := tgbotapi.NewUpdate(0)
	updCfg.AllowedUpdates = append(
		updCfg.AllowedUpdates,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	updCfg.Timeout = 60

	_, err = bot.Request(tgbotapi.NewSetMyCommands(botCommandRegister))
	if err != nil {
		return nil, errors.New(err.Error())
	}

	slog.Info(fmt.Sprintf("telegram:account:%s", bot.Self.UserName))

	slog.Info("telegram:init:finish")
	return &TelegramSource{
		sourceType: TelegramType,
		bot:        bot,
		router:     router,
		db:         db,
		cipher:     cipher,
		updChan:    bot.GetUpdatesChan(updCfg),
	}, nil
}

func (s *TelegramSource) GetSourceType() SourceType {
	return s.sourceType
}

func (s *TelegramSource) DoReq(ctx context.Context, wg *sync.WaitGroup) {
	slog.Info("telegram:do:req:start")
	defer wg.Done()
	defer slog.Info("telegram:do:req:finish")

	for {
		select {
		case upd := <-s.updChan:
			req := s.getReq(s.getMsg(&upd))
			if err := s.handleUpdate(ctx, req); err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				s.router.SendRes(req.ToResponse(false, err.Error()))
			} else {
				s.router.SendRes(req.ToResponse(true, ""))
			}
		case <-ctx.Done():
			slog.Info("telegram:do:req:ctx:done")
			return
		}
	}
}

func (s *TelegramSource) DoRes(ctx context.Context, wg *sync.WaitGroup) {
	slog.Info("telegram:do:res:start")
	defer wg.Done()
	defer slog.Info("telegram:do:res:finish")

	for {
		select {
		case res := <-s.router.ResChan():
			s.handleResponse(res)
		case <-ctx.Done():
			slog.Info("telegram:do:res:ctx:done")
			return
		}
	}
}

func (s *TelegramSource) getMsg(upd *tgbotapi.Update) *tgbotapi.Message {
	var message *tgbotapi.Message
	switch {
	case upd.Message != nil:
		message = upd.Message
	case upd.EditedMessage != nil:
		message = upd.EditedMessage
	}

	return message
}

func (s *TelegramSource) getReq(msg *tgbotapi.Message) *Request {
	if msg == nil {
		return nil
	}

	var content string
	var command string
	var reqType RequestType
	if msg.IsCommand() {
		content = msg.CommandArguments()
		command = msg.Command()
		reqType = RequestTypeCommand
	} else {
		content = msg.Text
		reqType = RequestTypeText
	}

	return &Request{
		TargetType:  JiraType,
		SourceType:  TelegramType,
		RequestType: reqType,
		Content:     content,
		Command:     command,
		UserId:      strconv.FormatInt(msg.From.ID, 10),
		ChatId:      strconv.FormatInt(msg.Chat.ID, 10),
		MessageId:   strconv.Itoa(msg.MessageID),
	}
}

func (s *TelegramSource) handleUpdate(ctx context.Context, req *Request) error {
	switch req.RequestType {
	case RequestTypeText:
		st, err := s.checkAccess(req)
		if err != nil {
			return err
		} else {
			switch st {
			case AccStateRegistered:
				s.router.SendReq(req)
			case AccStateNotRegistered:
				s.askToRegister(req)
			}
		}
	case RequestTypeCommand:
		switch req.Command {
		case botCommandRegister.Command:
			if err := s.handleRegistration(ctx, req); err != nil {
				return err
			}
		default:
			return errors.New("wrong command")
		}
	}
	return nil
}

func (s *TelegramSource) checkAccess(req *Request) (AccState, error) {
	ok, err := s.db.HasConnectionSource(s.sourceType.String(), req.UserId)
	if err != nil {
		return 0, err
	}
	if ok {
		return AccStateRegistered, nil
	} else {
		return AccStateNotRegistered, nil
	}
}

func (s *TelegramSource) handleRegistration(ctx context.Context, req *Request) error {
	ctxTx, cancelTx := context.WithCancel(ctx)
	defer cancelTx()

	slog.Info("reg:start", "source", s.sourceType.String(), "userId", req.UserId)

	tx, err := s.db.BeginTx(ctxTx)
	if err != nil {
		return err
	}

	accountGid, err := tx.CreateAccount()
	if err != nil {
		return err
	}

	err = tx.CreateSourceConnection(accountGid, s.sourceType.String(), req.UserId)
	if err != nil {
		return err
	}
	targetType := JiraType
	encoded, err := s.cipher.Encrypt(req.Content)
	if err != nil {
		return err
	}
	err = tx.CreateTargetConnection(accountGid, targetType.String(), encoded)
	if err != nil {
		return err
	}

	err = tx.CommitTx()
	if err != nil {
		return err
	}
	slog.Info("reg:finish", "source", s.sourceType.String(), "userId", req.UserId)

	return nil
}

func (s *TelegramSource) askToRegister(req *Request) {
	chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
	askMsg := tgbotapi.NewMessage(chatId, "/register <Private Access Token>")
	_, err = s.bot.Send(askMsg)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
}

func (s *TelegramSource) handleResponse(res *Response) {
	chatId, err := strconv.ParseInt(res.ChatId, 10, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
	var msg tgbotapi.MessageConfig
	if !res.Ok {
		msg = tgbotapi.NewMessage(chatId, "🙁")
	} else {
		msg = tgbotapi.NewMessage(chatId, "🙂")
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
}
