package internal

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"log/slog"
	"strconv"
	"wombat/pkg"
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
	cipher     *AesGcmCipher
	router     *Router
	db         *DbStorage
	updChan    tgbotapi.UpdatesChannel
}

var _ pkg.Task = (*TelegramSource)(nil)

func NewTelegramSource(token string, router *Router, db *DbStorage, cipher *AesGcmCipher) (ts *TelegramSource, err error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errors.WithStack(err)
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
		return nil, errors.WithStack(err)
	}

	slog.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

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

func (s *TelegramSource) Do(ctx context.Context) {
	for {
		select {
		case upd := <-s.updChan:
			msg := s.getMsg(&upd)
			req := s.getReq(msg)
			if !msg.IsCommand() {
				st, err := s.checkAccess(req)
				if err != nil {
					slog.Error(fmt.Sprintf("%+v", err))
					s.router.SendRes(req.ToResponse(false, err.Error()))
				} else {
					switch st {
					case Registered:
						s.router.SendReq(req)
					case NotRegistered:
						s.askToRegister(req)
					}
				}
			} else {
				switch msg.Command() {
				case botCommandRegister.Command:
					err := s.handleRegistration(ctx, req)
					if err != nil {
						slog.Error(fmt.Sprintf("%+v", err))
						s.router.SendRes(req.ToResponse(false, err.Error()))
					} else {
						s.router.SendRes(req.ToResponse(true, ""))
					}
				}
			}
		case res := <-s.router.ResChan():
			s.handleResponse(res)
		case <-ctx.Done():
			break
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
	if msg.IsCommand() {
		content = msg.CommandArguments()
	} else {
		content = msg.Text
	}

	return &Request{
		TargetType: JiraType,
		SourceType: TelegramType,
		Content:    content,
		UserId:     strconv.FormatInt(msg.From.ID, 10),
		ChatId:     strconv.FormatInt(msg.Chat.ID, 10),
		MessageId:  strconv.Itoa(msg.MessageID),
	}
}

func (s *TelegramSource) checkAccess(req *Request) (AccessState, error) {
	ok, err := s.db.HasConnectionSource(s.sourceType.String(), req.UserId)
	if err != nil {
		return 0, err
	}
	if ok {
		return Registered, nil
	} else {
		return NotRegistered, nil
	}
}

func (s *TelegramSource) handleRegistration(ctx context.Context, req *Request) error {
	ctxTx, cancelTx := context.WithCancel(ctx)
	defer cancelTx()

	slog.Info("REG:START", "source", s.sourceType.String(), "userId", req.UserId)

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
	slog.Info("REG:FINISH", "source", s.sourceType.String(), "userId", req.UserId)

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
