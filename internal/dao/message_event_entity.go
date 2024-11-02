package dao

import (
	"time"
	"wombat/internal/domain"
)

type MessageEventEntity struct {
	Hash       string    `db:"hash"`
	EventType  string    `db:"event_type"`
	SourceType string    `db:"source_type"`
	Text       string    `db:"text"`
	AuthorId   string    `db:"author_id"`
	ChatId     string    `db:"chat_id"`
	MessageId  string    `db:"message_id"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (receiver *MessageEventEntity) ToDomain() *domain.MessageEvent {
	var eventType domain.EventType
	eventType.FromString(receiver.EventType)
	var sourceType domain.SourceType
	sourceType.FromString(receiver.SourceType)
	return &domain.MessageEvent{
		Hash:       receiver.Hash,
		EventType:  eventType,
		SourceType: sourceType,
		Text:       receiver.Text,
		AuthorId:   receiver.AuthorId,
		ChatId:     receiver.ChatId,
		MessageId:  receiver.MessageId,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

func (receiver *MessageEventEntity) FromDomain(domain *domain.MessageEvent) {
	receiver.Hash = domain.Hash
	receiver.EventType = domain.EventType.String()
	receiver.SourceType = domain.SourceType.String()
	receiver.Text = domain.Text
	receiver.AuthorId = domain.AuthorId
	receiver.ChatId = domain.ChatId
	receiver.MessageId = domain.MessageId
	receiver.CreateTs = domain.CreateTs
	receiver.UpdateTs = domain.UpdateTs
}
