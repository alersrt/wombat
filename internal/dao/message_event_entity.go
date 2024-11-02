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
	return &domain.MessageEvent{
		Hash:       receiver.Hash,
		EventType:  domain.EventTypeFromString(receiver.EventType),
		SourceType: domain.SourceTypeFromString(receiver.SourceType),
		Text:       receiver.Text,
		AuthorId:   receiver.AuthorId,
		ChatId:     receiver.ChatId,
		MessageId:  receiver.MessageId,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

func MessageEventEntityFromDomain(domain *domain.MessageEvent) *MessageEventEntity {
	return &MessageEventEntity{
		Hash:       domain.Hash,
		EventType:  domain.EventType.String(),
		SourceType: domain.SourceType.String(),
		Text:       domain.Text,
		AuthorId:   domain.AuthorId,
		ChatId:     domain.ChatId,
		MessageId:  domain.MessageId,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
