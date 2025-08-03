package storage

import (
	"github.com/google/uuid"
	"time"
	"wombat/internal/domain"
)

type TargetConnectionEntity struct {
	Gid        uuid.UUID `db:"gid"`
	AccountGid uuid.UUID `db:"account_gid"`
	TargetType string    `db:"target_type"`
	Token      string    `db:"token"`
	CreateTs   time.Time `db:"create_ts"`
	UpdateTs   time.Time `db:"update_ts"`
}

func (receiver *TargetConnectionEntity) ToDomain() *domain.TargetConnection {
	return &domain.TargetConnection{
		Gid:        receiver.Gid,
		AccountGid: receiver.AccountGid,
		TargetType: domain.TargetTypeFromString(receiver.TargetType),
		Token:      receiver.Token,
		CreateTs:   receiver.CreateTs,
		UpdateTs:   receiver.UpdateTs,
	}
}

func (receiver *TargetConnectionEntity) FromDomain(domain *domain.TargetConnection) Entity[domain.TargetConnection] {
	return &TargetConnectionEntity{
		Gid:        domain.Gid,
		AccountGid: domain.AccountGid,
		TargetType: domain.TargetType.String(),
		Token:      domain.Token,
		CreateTs:   domain.CreateTs,
		UpdateTs:   domain.UpdateTs,
	}
}
