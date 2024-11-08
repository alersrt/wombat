package dao

import (
	"time"
	"wombat/internal/domain"
)

type CommentEntity struct {
	SourceType string    `db:"source_type"`
	Text       string    `db:"text"`
	AuthorId   string    `db:"author_id"`
	ChatId     string    `db:"chat_id"`
	MessageId  string    `db:"message_id"`
	CommentId  string    `db:"comment_id"`
	Tag        string    `db:"tag"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (receiver *CommentEntity) ToDomain() *domain.Comment {
	return &domain.Comment{
		Message: &domain.Message{
			SourceType: domain.SourceTypeFromString(receiver.SourceType),
			Text:       receiver.Text,
			AuthorId:   receiver.AuthorId,
			ChatId:     receiver.ChatId,
			MessageId:  receiver.MessageId,
		},
		Tag:       receiver.Tag,
		CommentId: receiver.CommentId,
		CreateTs:  receiver.CreateTs,
		UpdateTs:  receiver.UpdateTs,
	}
}

type CommentEntityFactory struct{}

func (receiver *CommentEntityFactory) EmptyEntity() Entity[domain.Comment] {
	return &CommentEntity{}
}

func (receiver *CommentEntityFactory) FromDomain(domain *domain.Comment) Entity[domain.Comment] {
	return &CommentEntity{
		SourceType: domain.SourceType.String(),
		Text:       domain.Text,
		AuthorId:   domain.AuthorId,
		ChatId:     domain.ChatId,
		MessageId:  domain.MessageId,
		CommentId:  domain.CommentId,
		Tag:        domain.Tag,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
