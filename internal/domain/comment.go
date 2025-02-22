package domain

import (
	"github.com/google/uuid"
	"time"
)

type Comment struct {
	Gid uuid.UUID `json:"gid"`
	*Message
	Tag       string    `json:"tag"`
	CommentId string    `json:"comment_id"`
	CreateTs  time.Time `json:"create_ts"`
	UpdateTs  time.Time `json:"update_ts"`
}
