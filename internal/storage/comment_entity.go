package storage

import (
	"github.com/google/uuid"
	"time"
	"wombat/internal/domain"
)

type CommentEntity struct {
	Gid        uuid.UUID `db:"gid"`
	SourceType string    `db:"source_type"`
	TargetType string    `db:"target_type"`
	UserId     string    `db:"user_id"`
	ChatId     string    `db:"chat_id"`
	MessageId  string    `db:"message_id"`
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
			UserId:     receiver.UserId,
			ChatId:     receiver.ChatId,
			MessageId:  receiver.MessageId,
			TargetType: domain.TargetTypeFromString(receiver.TargetType),
		},
		CommentId: receiver.CommentId,
		Tag:       receiver.Tag,
		CreateTs:  receiver.CreateTs,
		UpdateTs:  receiver.UpdateTs,
	}
}

func (receiver *CommentEntity) FromDomain(domain *domain.Comment) *CommentEntity {
	return &CommentEntity{
		Gid:        domain.Gid,
		TargetType: domain.TargetType.String(),
		SourceType: domain.SourceType.String(),
		CommentId:  domain.CommentId,
		UserId:     domain.UserId,
		ChatId:     domain.ChatId,
		MessageId:  domain.MessageId,
		Tag:        domain.Tag,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.CreateTs,
	}
}
