package domain

import (
	"github.com/google/uuid"
	"time"
)

type TargetConnection struct {
	Gid        uuid.UUID  `db:"gid"`
	AccountGid string     `db:"account_gid"`
	TargetType TargetType `db:"target_type"`
	Token      string     `db:"token"`
	CreateTs   time.Time  `db:"create_ts"`
	UpdateTs   time.Time  `db:"update_ts"`
}
