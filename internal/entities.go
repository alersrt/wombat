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

func (receiver *SourceConnectionEntity) ToDomain() *SourceConnection {
	return &SourceConnection{
		Gid:        receiver.Gid,
		AccountGid: receiver.AccountGid,
		SourceType: SourceTypeFromString(receiver.SourceType),
		UserId:     receiver.UserId,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

func (receiver *SourceConnectionEntity) FromDomain(domain *SourceConnection) Entity[SourceConnection] {
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
	Token      string    `db:"token"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (receiver *TargetConnectionEntity) ToDomain() *TargetConnection {
	return &TargetConnection{
		Gid:        receiver.Gid,
		AccountGid: receiver.AccountGid,
		TargetType: TargetTypeFromString(receiver.TargetType),
		Token:      receiver.Token,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

func (receiver *TargetConnectionEntity) FromDomain(domain *TargetConnection) Entity[TargetConnection] {
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

func (receiver *CommentEntity) ToDomain() *Comment {
	return &Comment{
		Gid: receiver.Gid,
		Message: &Message{
			SourceType: SourceTypeFromString(receiver.SourceType),
			UserId:     receiver.UserId,
			ChatId:     receiver.ChatId,
			MessageId:  receiver.MessageId,
			TargetType: TargetTypeFromString(receiver.TargetType),
		},
		CommentId: receiver.CommentId,
		Tag:       receiver.Tag,
		CreateTs:  receiver.CreateTs,
		UpdateTs:  receiver.UpdateTs,
	}
}

func (receiver *CommentEntity) FromDomain(domain *Comment) Entity[Comment] {
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
