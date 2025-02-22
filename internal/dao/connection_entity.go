package dao

import (
	"time"
	"wombat/internal/domain"
)

type ConnectionEntity struct {
	TargetType domain.TargetType `db:"target_type"`
	Token      string            `db:"token"`
	SourceType domain.SourceType `db:"source_type"`
	AuthorId   string            `db:"author_id"`
	CreateTs   time.Time         `db:"create_ts"`
	UpdateTs   time.Time         `db:"update_ts"`
}

func (receiver *ConnectionEntity) ToDomain() *domain.Connection {
	return &domain.Connection{
		TargetType: receiver.TargetType,
		Token:      receiver.Token,
		SourceType: receiver.SourceType,
		AuthorId:   receiver.AuthorId,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

type ConnectionEntityFactory struct{}

func (receiver *ConnectionEntityFactory) EmptyEntity() Entity[domain.Connection] {
	return &ConnectionEntity{}
}

func (receiver *ConnectionEntityFactory) FromDomain(domain *domain.Connection) Entity[domain.Connection] {
	return &ConnectionEntity{
		TargetType: domain.TargetType,
		Token:      domain.Token,
		SourceType: domain.SourceType,
		AuthorId:   domain.AuthorId,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
