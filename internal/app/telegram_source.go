package app

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log/slog"
	"strconv"
	"wombat/internal/domain"
)

const (
	COMMAND_REGISTER   = "register"
	COMMAND_UNREGISTER = "unregister"
)

type TelegramSource struct {
	*tgbotapi.BotAPI
	fwdChan chan *domain.Message
}

func NewTelegramSource(token string, fwdChan chan *domain.Message) (*TelegramSource, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	return &TelegramSource{
		BotAPI:  bot,
		fwdChan: fwdChan,
	}, nil
}

func (receiver *TelegramSource) GetSourceType() domain.SourceType {
	return domain.TELEGRAM
}

func (receiver *TelegramSource) Init() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.AllowedUpdates = append(
		u.AllowedUpdates,
		tgbotapi.UpdateTypeMessage,
		tgbotapi.UpdateTypeEditedMessage,
	)
	u.Timeout = 60

	commands := tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     COMMAND_REGISTER,
			Description: COMMAND_REGISTER,
		}, tgbotapi.BotCommand{
			Command:     COMMAND_UNREGISTER,
			Description: COMMAND_UNREGISTER,
		},
	)
	_, err := receiver.Request(commands)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	slog.Info(fmt.Sprintf("Authorized on account %s", receiver.Self.UserName))

	return receiver.GetUpdatesChan(u), nil
}

func (receiver *TelegramSource) Process() {
	updates, err := receiver.Init()
	if err != nil {
		slog.Error(err.Error())
		return
	}

	for update := range updates {

		var message *tgbotapi.Message
		switch {
		case update.Message != nil:
			message = update.Message
		case update.EditedMessage != nil:
			message = update.EditedMessage
		}

		if !update.Message.IsCommand() {
			receiver.handleMessage(message)
		} else {
			switch message.Command() {
			case COMMAND_REGISTER:
				receiver.handleRegistration(message.CommandArguments())
			}
		}
	}
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

func (receiver *TelegramSource) handleRegistration(arg string) {
	slog.Info("Example", "token", arg)
}
