package domain

import (
	"github.com/google/uuid"
	"time"
)

type TargetConnection struct {
	Gid        uuid.UUID  `json:"gid"`
	AccountGid uuid.UUID  `json:"account_gid"`
	TargetType TargetType `json:"target_type"`
	Token      string     `json:"token"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
