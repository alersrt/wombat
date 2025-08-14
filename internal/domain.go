package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type SourceType uint

type TargetType uint

const (
	TelegramType SourceType = iota
)

const (
	JiraType TargetType = iota
)

var (
	sourceTypeToString = map[SourceType]string{
		TelegramType: "TELEGRAM",
	}
	sourceTypeFromString = map[string]SourceType{
		"TELEGRAM": TelegramType,
	}
	targetTypeToString = map[TargetType]string{
		JiraType: "JIRA",
	}
	targetTypeFromString = map[string]TargetType{
		"JIRA": JiraType,
	}
)

func (t *SourceType) String() string {
	return sourceTypeToString[*t]
}

func SourceTypeFromString(s string) SourceType {
	return sourceTypeFromString[s]
}

func (t *SourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(sourceTypeToString[*t])
}

func (t *SourceType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = sourceTypeFromString[s]
	return nil
}

func (t *TargetType) String() string {
	return targetTypeToString[*t]
}

func TargetTypeFromString(s string) TargetType {
	return targetTypeFromString[s]
}

func (t *TargetType) MarshalJSON() ([]byte, error) {
	return json.Marshal(targetTypeToString[*t])
}

func (t *TargetType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*t = targetTypeFromString[s]
	return nil
}

type AccessState uint

const (
	NotRegistered AccessState = iota
	Registered
)

type Request struct {
	SourceType SourceType `json:"source_type"`
	TargetType TargetType `json:"target_type"`
	Content    string     `json:"content"`
	UserId     string     `json:"user_id"`
	ChatId     string     `json:"chat_id"`
	MessageId  string     `json:"message_id"`
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
