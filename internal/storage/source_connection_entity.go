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

func (receiver *SourceConnectionEntity) ToDomain() *domain.SourceConnection {
	return &domain.SourceConnection{
		Gid:        receiver.Gid,
		AccountGid: receiver.AccountGid,
		SourceType: domain.SourceTypeFromString(receiver.SourceType),
		UserId:     receiver.UserId,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

func (receiver *SourceConnectionEntity) FromDomain(domain *domain.SourceConnection) *SourceConnectionEntity {
	return &SourceConnectionEntity{
		Gid:        domain.Gid,
		AccountGid: domain.AccountGid,
		SourceType: domain.SourceType.String(),
		UserId:     domain.UserId,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
