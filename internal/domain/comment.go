package domain

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Gid uuid.UUID `json:"gid"`
	*Message
	Tag        string     `json:"tag"`
	TargetType TargetType `json:"target_type"`
	CommentId  string     `json:"comment_id"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
