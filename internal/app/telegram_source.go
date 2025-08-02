package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"log/slog"
	"strconv"
	"wombat/internal/domain"
	"wombat/internal/storage"
)

var (
	BotCommandRegister = tgbotapi.BotCommand{
		Command:     "register",
		Description: "RG",
	}
)

type TelegramSource struct {
	sourceType domain.SourceType
	*tgbotapi.BotAPI
	fwdChan chan *domain.Message
	db      *storage.DbStorage
}

func NewTelegramSource(token string, fwdChan chan *domain.Message, db *storage.DbStorage) (*TelegramSource, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &TelegramSource{
		sourceType: domain.TELEGRAM,
		BotAPI:     bot,
		fwdChan:    fwdChan,
		db:         db,
	}, nil
}

func (receiver *TelegramSource) GetSourceType() domain.SourceType {
	return receiver.sourceType
}

func (receiver *TelegramSource) init() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = append(
		u.AllowedUpdates,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	u.Timeout = 60

	commands := tgbotapi.NewSetMyCommands(BotCommandRegister)
	_, err := receiver.Request(commands)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info(fmt.Sprintf("Authorized on account %s", receiver.Self.UserName))

	return receiver.GetUpdatesChan(u), nil
}

func (receiver *TelegramSource) Process() {
	updates, err := receiver.init()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for update := range updates {

		message := receiver.getMessage(&update)
		err := receiver.checkAccess(message)
		if err != nil {
			slog.Warn(err.Error())
			return
		}

		if !update.Message.IsCommand() {
			receiver.handleMessage(message)
		} else {
			switch message.Command() {
			case BotCommandRegister.Command:
				receiver.handleRegistration(message.CommandArguments())
			}
		}
	}
}

func (receiver *TelegramSource) checkAccess(message *tgbotapi.Message) error {
	isOk := receiver.db.HasConnectionSource(receiver.sourceType.String(), strconv.FormatInt(message.From.ID, 10))
	if isOk {
		return nil
	} else {
		return errors.Errorf(domain.ErrorNotRegistered)
	}
}

func (receiver *TelegramSource) getMessage(update *tgbotapi.Update) *tgbotapi.Message {
	var message *tgbotapi.Message
	switch {
	case update.Message != nil:
		message = update.Message
	case update.EditedMessage != nil:
		message = update.EditedMessage
	}
	return message
}

func (receiver *TelegramSource) handleMessage(message *tgbotapi.Message) {
	msg := &domain.Message{
		TargetType: domain.JIRA,
		SourceType: domain.TELEGRAM,
		Content:    message.Text,
		UserId:     strconv.FormatInt(message.From.ID, 10),
		ChatId:     strconv.FormatInt(message.Chat.ID, 10),
		MessageId:  strconv.Itoa(message.MessageID),
	}
	receiver.fwdChan <- msg
}

func (receiver *TelegramSource) handleRegistration(token string) {
	slog.Info("Example", "token", token)
}
