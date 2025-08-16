package domain

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type SourceType string

type TargetType string

const (
	SourceTypeTelegram = "TELEGRAM"
	TargetTypeJira     = "JIRA"
)

func (t *SourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *SourceType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("source_type:unmarchal: %w", err)
	}
	*t = SourceType(s)
	return nil
}

func (t *TargetType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TargetType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return fmt.Errorf("target_type:unmarchal: %w", err)
	}
	*t = TargetType(s)
	return nil
}

type AccState string

const (
	AccStateNotRegistered = "NOT_REGISTERED"
	AccStateRegistered    = "REGISTERED"
)

type RequestType string

const (
	RequestTypeText    = "TEXT"
	RequestTypeCommand = "COMMAND"
)

type Request struct {
	SourceType  SourceType  `json:"source_type"`
	TargetType  TargetType  `json:"target_type"`
	RequestType RequestType `json:"request_type"`
	Content     string      `json:"content"`
	Command     string      `json:"command"`
	UserId      string      `json:"user_id"`
	ChatId      string      `json:"chat_id"`
	MessageId   string      `json:"message_id"`
}

func (r *Request) ToResponse(ok bool, desc string) *Response {
	return &Response{
		Ok:         ok,
		Desc:       desc,
		SourceType: r.SourceType,
		TargetType: r.TargetType,
		UserId:     r.UserId,
		ChatId:     r.ChatId,
		MessageId:  r.MessageId,
	}
}

type Response struct {
	Ok         bool       `json:"ok"`
	Desc       string     `json:"desc"`
	SourceType SourceType `json:"source_type"`
	TargetType TargetType `json:"target_type"`
	UserId     string     `json:"user_id"`
	ChatId     string     `json:"chat_id"`
	MessageId  string     `json:"message_id"`
}

type Comment struct {
	Gid uuid.UUID `json:"gid"`
	*Request
	Tag       string    `json:"tag"`
	CommentId string    `json:"comment_id"`
	CreateTs  time.Time `json:"create_ts"`
	UpdateTs  time.Time `json:"update_ts"`
}

type SourceConnection struct {
	Gid        uuid.UUID  `json:"gid"`
	AccountGid uuid.UUID  `json:"account_gid"`
	SourceType SourceType `json:"source_type"`
	UserId     string     `json:"user_id"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}

type TargetConnection struct {
	Gid        uuid.UUID  `json:"gid"`
	AccountGid uuid.UUID  `json:"account_gid"`
	TargetType TargetType `json:"target_type"`
	Token      []byte     `json:"token"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
