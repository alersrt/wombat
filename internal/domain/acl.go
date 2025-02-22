package domain

import (
	"github.com/google/uuid"
	"time"
)

type Acl struct {
	Gid        uuid.UUID  `json:"gid"`
	SourceType SourceType `json:"source_type"`
	AuthorId   string     `json:"author_id"`
	IsAllowed  bool       `json:"is_allowed"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
