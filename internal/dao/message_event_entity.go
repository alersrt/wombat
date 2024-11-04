package dao

import (
	"time"
	"wombat/internal/domain"
)

type MessageEventEntity struct {
	Hash       string    `db:"hash"`
	SourceType string    `db:"source_type"`
	Text       string    `db:"text"`
	AuthorId   string    `db:"author_id"`
	ChatId     string    `db:"chat_id"`
	MessageId  string    `db:"message_id"`
	CommentId  string    `db:"comment_id"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (receiver *MessageEventEntity) ToDomain() *domain.MessageEvent {
	return &domain.MessageEvent{
		Hash:       receiver.Hash,
		SourceType: domain.SourceTypeFromString(receiver.SourceType),
		Text:       receiver.Text,
		AuthorId:   receiver.AuthorId,
		ChatId:     receiver.ChatId,
		MessageId:  receiver.MessageId,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
		CommentId:  receiver.CommentId,
	}
}

type MessageEventEntityFactory struct{}

func (receiver *MessageEventEntityFactory) EmptyEntity() Entity[domain.MessageEvent] {
	return &MessageEventEntity{}
}

func (receiver *MessageEventEntityFactory) FromDomain(domain *domain.MessageEvent) Entity[domain.MessageEvent] {
	return &MessageEventEntity{
		Hash:       domain.Hash,
		SourceType: domain.SourceType.String(),
		Text:       domain.Text,
		AuthorId:   domain.AuthorId,
		ChatId:     domain.ChatId,
		MessageId:  domain.MessageId,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
		CommentId:  domain.CommentId,
	}
}
