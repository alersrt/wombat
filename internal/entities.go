package internal

import (
	"github.com/google/uuid"
	"time"
)

type SourceConnectionEntity struct {
	Gid        uuid.UUID `db:"gid"`
	AccountGid uuid.UUID `db:"account_gid"`
	SourceType string    `db:"source_type"`
	UserId     string    `db:"user_id"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (e *SourceConnectionEntity) ToDomain() *SourceConnection {
	return &SourceConnection{
		Gid:        e.Gid,
		AccountGid: e.AccountGid,
		SourceType: SourceTypeFromString(e.SourceType),
		UserId:     e.UserId,
		CreateTs:   e.CreateTs,
		UpdateTs:   e.UpdateTs,
	}
}

func (e *SourceConnectionEntity) FromDomain(domain *SourceConnection) Entity[SourceConnection] {
	return &SourceConnectionEntity{
		Gid:        domain.Gid,
		AccountGid: domain.AccountGid,
		SourceType: domain.SourceType.String(),
		UserId:     domain.UserId,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}

type TargetConnectionEntity struct {
	Gid        uuid.UUID `db:"gid"`
	AccountGid uuid.UUID `db:"account_gid"`
	TargetType string    `db:"target_type"`
	Token      []byte    `db:"token"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (e *TargetConnectionEntity) ToDomain() *TargetConnection {
	return &TargetConnection{
		Gid:        e.Gid,
		AccountGid: e.AccountGid,
		TargetType: TargetTypeFromString(e.TargetType),
		Token:      e.Token,
		CreateTs:   e.CreateTs,
		UpdateTs:   e.UpdateTs,
	}
}

func (e *TargetConnectionEntity) FromDomain(domain *TargetConnection) Entity[TargetConnection] {
	return &TargetConnectionEntity{
		Gid:        domain.Gid,
		AccountGid: domain.AccountGid,
		TargetType: domain.TargetType.String(),
		Token:      domain.Token,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}

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

func (e *CommentEntity) ToDomain() *Comment {
	return &Comment{
		Gid: e.Gid,
		Message: &Message{
			SourceType: SourceTypeFromString(e.SourceType),
			UserId:     e.UserId,
			ChatId:     e.ChatId,
			MessageId:  e.MessageId,
			TargetType: TargetTypeFromString(e.TargetType),
		},
		CommentId: e.CommentId,
		Tag:       e.Tag,
		CreateTs:  e.CreateTs,
		UpdateTs:  e.UpdateTs,
	}
}

func (e *CommentEntity) FromDomain(domain *Comment) Entity[Comment] {
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
