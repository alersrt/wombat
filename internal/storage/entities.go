package storage

import (
	"github.com/google/uuid"
	"time"
	"wombat/internal/domain"
)

type SourceConnectionEntity struct {
	Gid        uuid.UUID `db:"gid"`
	AccountGid uuid.UUID `db:"account_gid"`
	SourceType string    `db:"source_type"`
	UserId     string    `db:"user_id"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (e *SourceConnectionEntity) ToDomain() *domain.SourceConnection {
	return &domain.SourceConnection{
		Gid:        e.Gid,
		AccountGid: e.AccountGid,
		SourceType: domain.SourceType(e.SourceType),
		UserId:     e.UserId,
		CreateTs:   e.CreateTs,
		UpdateTs:   e.UpdateTs,
	}
}

func (e *SourceConnectionEntity) FromDomain(domain *domain.SourceConnection) Entity[domain.SourceConnection] {
	return &SourceConnectionEntity{
		Gid:        domain.Gid,
		AccountGid: domain.AccountGid,
		SourceType: string(domain.SourceType),
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

func (e *TargetConnectionEntity) ToDomain() *domain.TargetConnection {
	return &domain.TargetConnection{
		Gid:        e.Gid,
		AccountGid: e.AccountGid,
		TargetType: domain.TargetType(e.TargetType),
		Token:      e.Token,
		CreateTs:   e.CreateTs,
		UpdateTs:   e.UpdateTs,
	}
}

func (e *TargetConnectionEntity) FromDomain(domain *domain.TargetConnection) Entity[domain.TargetConnection] {
	return &TargetConnectionEntity{
		Gid:        domain.Gid,
		AccountGid: domain.AccountGid,
		TargetType: string(domain.TargetType),
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

func (e *CommentEntity) ToDomain() *domain.Comment {
	return &domain.Comment{
		Gid: e.Gid,
		Request: &domain.Request{
			SourceType: domain.SourceType(e.SourceType),
			TargetType: domain.TargetType(e.TargetType),
			UserId:     e.UserId,
			ChatId:     e.ChatId,
			MessageId:  e.MessageId,
		},
		CommentId: e.CommentId,
		Tag:       e.Tag,
		CreateTs:  e.CreateTs,
		UpdateTs:  e.UpdateTs,
	}
}

func (e *CommentEntity) FromDomain(domain *domain.Comment) Entity[domain.Comment] {
	return &CommentEntity{
		Gid:        domain.Gid,
		TargetType: string(domain.TargetType),
		SourceType: string(domain.SourceType),
		CommentId:  domain.CommentId,
		UserId:     domain.UserId,
		ChatId:     domain.ChatId,
		MessageId:  domain.MessageId,
		Tag:        domain.Tag,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.CreateTs,
	}
}
