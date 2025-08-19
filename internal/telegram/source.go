package telegram

import (
	"context"
	"fmt"
	api "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strconv"
	"wombat/internal/domain"
	"wombat/internal/router"
	"wombat/internal/storage"
	"wombat/pkg/cipher"
)

var (
	botCommandRegister = api.BotCommand{Command: "register", Description: "RG"}
)

type Source struct {
	sourceType domain.SourceType
	cipher     *cipher.AesGcmCipher
	rt         *router.Router
	db         *storage.DbStorage
	bot        *api.BotAPI
	cfg        api.UpdateConfig
}

func NewTelegramSource(
	token string,
	rt *router.Router,
	db *storage.DbStorage,
	cipher *cipher.AesGcmCipher,
) (*Source, error) {
	slog.Info("tg:new:start")

	bot, err := api.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("tg:new:err: %w", err)
	}

	updCfg := api.NewUpdate(0)
	updCfg.AllowedUpdates = append(
		updCfg.AllowedUpdates,
		api.UpdateTypeMessage,
		api.UpdateTypeEditedMessage,
	)
	updCfg.Timeout = 60

	_, err = bot.Request(api.NewSetMyCommands(botCommandRegister))
	if err != nil {
		return nil, fmt.Errorf("tg:new:err: %w", err)
	}

	slog.Info(fmt.Sprintf("tg:new:acc: %s", bot.Self.UserName))

	slog.Info("tg:new:finish")
	return &Source{
		sourceType: domain.SourceTypeTelegram,
		bot:        bot,
		rt:         rt,
		db:         db,
		cipher:     cipher,
		cfg:        updCfg,
	}, nil
}

func (s *Source) DoReq(ctx context.Context) {
	slog.Info("tg:do:req:start")
	defer slog.Info("tg:do:req:finish")

	for {
		select {
		case upd := <-s.bot.GetUpdatesChan(s.cfg):
			req := s.updToReq(&upd)
			if err := s.handleUpdate(ctx, req); err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				s.rt.SendRes(req.ToResponse(false, err.Error()))
			} else {
				s.rt.SendRes(req.ToResponse(true, ""))
			}
		case <-ctx.Done():
			s.bot.StopReceivingUpdates()
			slog.Info("tg:do:req:ctx:done")
			return
		}
	}
}

func (s *Source) DoRes(ctx context.Context) {
	slog.Info("tg:do:res:start")
	defer slog.Info("tg:do:res:finish")

	for {
		select {
		case res := <-s.rt.ResChan():
			s.handleResponse(res)
		case <-ctx.Done():
			slog.Info("tg:do:res:ctx:done")
			return
		}
	}
}

func (s *Source) updToReq(upd *api.Update) *domain.Request {
	var msg *api.Message
	switch {
	case upd.Message != nil:
		msg = upd.Message
	case upd.EditedMessage != nil:
		msg = upd.EditedMessage
	}

	if msg == nil {
		return nil
	}

	var content string
	var command string
	var reqType domain.RequestType
	if msg.IsCommand() {
		content = msg.CommandArguments()
		command = msg.Command()
		reqType = domain.RequestTypeCommand
	} else {
		content = msg.Text
		reqType = domain.RequestTypeText
	}

	return &domain.Request{
		TargetType:  domain.TargetTypeJira,
		SourceType:  domain.SourceTypeTelegram,
		RequestType: reqType,
		Content:     content,
		Command:     command,
		UserId:      strconv.FormatInt(msg.From.ID, 10),
		ChatId:      strconv.FormatInt(msg.Chat.ID, 10),
		MessageId:   strconv.Itoa(msg.MessageID),
	}
}

func (s *Source) handleUpdate(ctx context.Context, req *domain.Request) error {
	switch req.RequestType {
	case domain.RequestTypeText:
		st, err := s.checkAccess(req)
		if err != nil {
			return err
		} else {
			switch st {
			case domain.AccStateRegistered:
				s.rt.SendReq(req)
			case domain.AccStateNotRegistered:
				s.askToRegister(req)
			}
		}
	case domain.RequestTypeCommand:
		switch req.Command {
		case botCommandRegister.Command:
			if err := s.handleRegistration(ctx, req); err != nil {
				return err
			}
		default:
			return fmt.Errorf("tg:upd: wrong cmd=%s", req.Command)
		}
	}
	return nil
}

func (s *Source) checkAccess(req *domain.Request) (domain.AccState, error) {
	ok, err := s.db.HasConnectionSource(string(s.sourceType), req.UserId)
	if err != nil {
		return "", err
	}
	if ok {
		return domain.AccStateRegistered, nil
	} else {
		return domain.AccStateNotRegistered, nil
	}
}

func (s *Source) handleRegistration(ctx context.Context, req *domain.Request) error {
	ctxTx, cancelTx := context.WithCancel(ctx)
	defer cancelTx()

	slog.Info("tg:reg:start", "source", s.sourceType, "userId", req.UserId)

	tx, err := s.db.BeginTx(ctxTx)
	if err != nil {
		return err
	}

	accountGid, err := tx.CreateAccount()
	if err != nil {
		return err
	}

	err = tx.CreateSourceConnection(accountGid, string(s.sourceType), req.UserId)
	if err != nil {
		return err
	}
	targetType := domain.TargetTypeJira
	encoded, err := s.cipher.Encrypt(req.Content)
	if err != nil {
		return err
	}
	err = tx.CreateTargetConnection(accountGid, targetType, encoded)
	if err != nil {
		return err
	}

	err = tx.CommitTx()
	if err != nil {
		return err
	}
	slog.Info("tg:reg:finish", "source", s.sourceType, "userId", req.UserId)

	return nil
}

func (s *Source) askToRegister(req *domain.Request) {
	chatId, err := strconv.ParseInt(req.ChatId, 10, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
	askMsg := api.NewMessage(chatId, "/register <Private Access Token>")
	_, err = s.bot.Send(askMsg)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
	}
}

func (s *Source) handleResponse(res *domain.Response) {
	chatId, err := strconv.ParseInt(res.ChatId, 10, 64)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		return
	}
	var msg api.MessageConfig
	if !res.Ok {
		msg = api.NewMessage(chatId, "🙁")
	} else {
		msg = api.NewMessage(chatId, "🙂")
	}
	_, err = s.bot.Send(msg)
	if err != nil {
		slog.Error(fmt.Sprintf("%+v", err))
		return
	}
}
