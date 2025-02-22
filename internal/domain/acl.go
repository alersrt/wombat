package domain

import "time"

type Acl struct {
	SourceType SourceType `json:"source_type"`
	AuthorId   string     `json:"author_id"`
	IsAllowed  bool       `json:"is_allowed"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
