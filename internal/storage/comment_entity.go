package storage

import (
	"github.com/google/uuid"
	"time"
	"wombat/internal/domain"
)

type CommentEntity struct {
	Gid        uuid.UUID `db:"gid"`
	SourceType string    `db:"source_type"`
	Text       string    `db:"text"`
	AuthorId   string    `db:"author_id"`
	ChatId     string    `db:"chat_id"`
	MessageId  string    `db:"message_id"`
	TargetType string    `db:"target_type"`
	CommentId  string    `db:"comment_id"`
	Tag        string    `db:"tag"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (receiver *CommentEntity) ToDomain() *domain.Comment {
	return &domain.Comment{
		Gid: receiver.Gid,
		Message: &domain.Message{
			SourceType: domain.SourceTypeFromString(receiver.SourceType),
			Text:       receiver.Text,
			AuthorId:   receiver.AuthorId,
			ChatId:     receiver.ChatId,
			MessageId:  receiver.MessageId,
			TargetType: domain.TargetTypeFromString(receiver.TargetType),
		},
		Tag:       receiver.Tag,
		CommentId: receiver.CommentId,
		CreateTs:  receiver.CreateTs,
		UpdateTs:  receiver.UpdateTs,
	}
}

func (receiver *CommentEntity) FromDomain(domain *domain.Comment) Entity[domain.Comment] {
	return &CommentEntity{
		Gid:        domain.Gid,
		SourceType: domain.SourceType.String(),
		Text:       domain.Text,
		AuthorId:   domain.AuthorId,
		ChatId:     domain.ChatId,
		MessageId:  domain.MessageId,
		TargetType: domain.TargetType.String(),
		CommentId:  domain.CommentId,
		Tag:        domain.Tag,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
