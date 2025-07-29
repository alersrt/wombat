package storage

import (
	"github.com/google/uuid"
	"time"
	"wombat/internal/domain"
)

type ConnectionEntity struct {
	Gid        uuid.UUID `db:"gid"`
	TargetType string    `db:"target_type"`
	Token      string    `db:"token"`
	SourceType string    `db:"source_type"`
	AuthorId   string    `db:"author_id"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (receiver *ConnectionEntity) ToDomain() *domain.Connection {
	return &domain.Connection{
		Gid:        receiver.Gid,
		TargetType: domain.TargetTypeFromString(receiver.TargetType),
		Token:      receiver.Token,
		SourceType: domain.SourceTypeFromString(receiver.SourceType),
		AuthorId:   receiver.AuthorId,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

func (receiver *ConnectionEntity) FromDomain(domain *domain.Connection) Entity[domain.Connection] {
	return &ConnectionEntity{
		Gid:        domain.Gid,
		TargetType: domain.TargetType.String(),
		Token:      domain.Token,
		SourceType: domain.SourceType.String(),
		AuthorId:   domain.AuthorId,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
