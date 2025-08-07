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
	fwdChan chan *Message
	db      *DbStorage
}

func NewTelegramSource(token string, fwdChan chan *Message, db *DbStorage, cipher *AesGcmCipher) (ts *TelegramSource, err error) {
	defer pkg.CatchWithReturn(&err)

	bot, err := tgbotapi.NewBotAPI(token)
	pkg.Throw(err)

	return &TelegramSource{
		sourceType: TelegramType,
		BotAPI:     bot,
		fwdChan:    fwdChan,
		db:         db,
		cipher:     cipher,
	}, nil
}

func (s *TelegramSource) GetSourceType() SourceType {
	return s.sourceType
}

func (s *TelegramSource) Do(ctx context.Context) (err error) {
	updates := s.init()
	defer pkg.CatchWithReturn(&err)

	for update := range updates {
		msg := s.getMessage(&update)

		if !update.Message.IsCommand() {
			switch s.checkAccess(msg) {
			case Registered:
				s.handleMessage(msg)
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
	}
	return
}

func (s *TelegramSource) init() tgbotapi.UpdatesChannel {
	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = append(
		u.AllowedUpdates,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	u.Timeout = 60

	commands := tgbotapi.NewSetMyCommands(botCommandRegister)
	_, err := s.Request(commands)
	pkg.Throw(err)

	slog.Info(fmt.Sprintf("Authorized on account %s", s.Self.UserName))

	return s.GetUpdatesChan(u)
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

func (s *TelegramSource) handleMessage(message *tgbotapi.Message) {
	msg := &Message{
		TargetType: JiraType,
		SourceType: TelegramType,
		Content:    message.Text,
		UserId:     strconv.FormatInt(message.From.ID, 10),
		ChatId:     strconv.FormatInt(message.Chat.ID, 10),
		MessageId:  strconv.Itoa(message.MessageID),
	}
	s.fwdChan <- msg
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
