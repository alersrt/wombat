package internal

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	*tgbotapi.BotAPI
	cipher  *AesGcmCipher
	router  *Router
	db      *DbStorage
	updChan tgbotapi.UpdatesChannel
}

func NewTelegramSource(token string, router *Router, db *DbStorage, cipher *AesGcmCipher) (ts *TelegramSource, err error) {
	defer pkg.CatchWithReturn(&err)

	bot, err := tgbotapi.NewBotAPI(token)
	pkg.Throw(err)

	updCfg := tgbotapi.NewUpdate(0)
	updCfg.AllowedUpdates = append(
		updCfg.AllowedUpdates,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	updCfg.Timeout = 60

	_, err = bot.Request(tgbotapi.NewSetMyCommands(botCommandRegister))
	pkg.Throw(err)

	slog.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	return &TelegramSource{
		sourceType: TelegramType,
		BotAPI:     bot,
		router:     router,
		db:         db,
		cipher:     cipher,
		updChan:    bot.GetUpdatesChan(updCfg),
	}, nil
}

func (s *TelegramSource) GetSourceType() SourceType {
	return s.sourceType
}

func (s *TelegramSource) Do(ctx context.Context) (err error) {
	defer pkg.CatchWithReturn(&err)

	select {
	case upd := <-s.updChan:
		msg := s.getMessage(&upd)
		if !upd.Message.IsCommand() {
			switch s.checkAccess(msg) {
			case Registered:
				s.handleRequest(msg)
			case NotRegistered:
				s.askToRegister(msg)
			}
		} else {
			switch msg.Command() {
			case botCommandRegister.Command:
				err := s.handleRegistration(ctx, strconv.FormatInt(msg.From.ID, 10), msg.CommandArguments())
				pkg.Throw(err)
			}
		}
	case res := <-s.router.GetRes():
		s.handleResponse(res)
	case <-ctx.Done():
		pkg.Throw(ctx.Err())
	}
	return
}

func (s *TelegramSource) checkAccess(message *tgbotapi.Message) AccessState {
	isOk := s.db.HasConnectionSource(s.sourceType.String(), strconv.FormatInt(message.From.ID, 10))
	if isOk {
		return Registered
	} else {
		return NotRegistered
	}
}

func (s *TelegramSource) getMessage(update *tgbotapi.Update) *tgbotapi.Message {
	var message *tgbotapi.Message
	switch {
	case update.Message != nil:
		message = update.Message
	case update.EditedMessage != nil:
		message = update.EditedMessage
	}
	return message
}

func (s *TelegramSource) handleRequest(message *tgbotapi.Message) {
	req := &Request{
		TargetType: JiraType,
		SourceType: TelegramType,
		Content:    message.Text,
		UserId:     strconv.FormatInt(message.From.ID, 10),
		ChatId:     strconv.FormatInt(message.Chat.ID, 10),
		MessageId:  strconv.Itoa(message.MessageID),
	}
	s.router.SendReq(req)
}

func (s *TelegramSource) handleResponse(res *Response) {
	chatId, err := strconv.ParseInt(res.ChatId, 10, 64)
	pkg.Throw(err)
	var msg tgbotapi.MessageConfig
	if res.Ok {
		msg = tgbotapi.NewMessage(chatId, "✅")
	} else {
		msg = tgbotapi.NewMessage(chatId, "❌")
	}
	_, err = s.Send(msg)
	pkg.Throw(err)
}

func (s *TelegramSource) askToRegister(message *tgbotapi.Message) {
	askMsg := tgbotapi.NewMessage(message.Chat.ID, "/register <Private Access Token>")
	_, err := s.Send(askMsg)
	pkg.Throw(err)
}

func (s *TelegramSource) handleRegistration(ctx context.Context, userId string, token string) (err error) {
	slog.Info("REG:START", "source", s.sourceType.String(), "userId", userId)

	tx := s.db.BeginTx(ctx)
	defer pkg.CatchWithReturnAndCall(&err, tx.RollbackTx)

	accountGid := tx.CreateAccount()
	tx.CreateSourceConnection(accountGid, s.sourceType.String(), userId)
	targetType := JiraType
	tx.CreateTargetConnection(accountGid, targetType.String(), s.cipher.Encrypt(token))

	tx.CommitTx()
	slog.Info("REG:FINISH", "source", s.sourceType.String(), "userId", userId)
	return
}
