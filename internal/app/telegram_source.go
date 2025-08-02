package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		sourceType: domain.Telegram,
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

		if !update.Message.IsCommand() {
			switch receiver.checkAccess(message) {
			case domain.Registered:
				receiver.handleMessage(message)
			case domain.NotRegistered:
				receiver.askToRegister(message)
			}
		} else {
			switch message.Command() {
			case BotCommandRegister.Command:
				receiver.handleRegistration(strconv.FormatInt(message.From.ID, 10), message.CommandArguments())
			}
		}
	}
}

func (receiver *TelegramSource) checkAccess(message *tgbotapi.Message) domain.AccessState {
	isOk := receiver.db.HasConnectionSource(receiver.sourceType.String(), strconv.FormatInt(message.From.ID, 10))
	if isOk {
		return domain.Registered
	} else {
		return domain.NotRegistered
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
		TargetType: domain.Jira,
		SourceType: domain.Telegram,
		Content:    message.Text,
		UserId:     strconv.FormatInt(message.From.ID, 10),
		ChatId:     strconv.FormatInt(message.Chat.ID, 10),
		MessageId:  strconv.Itoa(message.MessageID),
	}
	receiver.fwdChan <- msg
}

func (receiver *TelegramSource) askToRegister(message *tgbotapi.Message) {
	askMsg := tgbotapi.NewMessage(message.Chat.ID, "/register <Private Access Token>")
	_, err := receiver.Send(askMsg)
	if err != nil {
		slog.Warn(err.Error())
		return
	}
}

func (receiver *TelegramSource) handleRegistration(userId string, token string) error {
	slog.Info("REG:START", "source", receiver.sourceType.String(), "userId", userId)

	tx, err := receiver.db.BeginTx()
	if err != nil {
		return err
	}

	accountGid, err := tx.CreateAccount()
	if err != nil {
		return tx.RollbackTx()
	}
	err = tx.CreateSourceConnection(accountGid, receiver.sourceType.String(), userId)
	if err != nil {
		return tx.RollbackTx()
	}
	targetType := domain.Jira
	err = tx.CreateTargetConnection(accountGid, targetType.String(), token)

	if err != nil {
		return tx.RollbackTx()
	}

	slog.Info("REG:FINISH", "source", receiver.sourceType.String(), "userId", userId)
	return tx.CommitTx()
}
