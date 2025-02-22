package domain

import "time"

type Connection struct {
	TargetType TargetType `json:"target_type"`
	Token      string     `json:"token"`
	SourceType SourceType `json:"source_type"`
	AuthorId   string     `json:"author_id"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
