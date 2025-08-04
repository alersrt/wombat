package internal

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type SourceType uint

type Source interface {
	GetSourceType() SourceType
	Process()
}

const (
	TelegramType SourceType = iota
)

var (
	sourceTypeToString = map[SourceType]string{
		TelegramType: "TELEGRAM",
	}
	sourceTypeFromString = map[string]SourceType{
		"TELEGRAM": TelegramType,
	}
)

func (receiver *SourceType) String() string {
	return sourceTypeToString[*receiver]
}

func SourceTypeFromString(s string) SourceType {
	return sourceTypeFromString[s]
}

func (receiver *SourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(sourceTypeToString[*receiver])
}

func (receiver *SourceType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*receiver = sourceTypeFromString[s]
	return nil
}

type Target interface {
	GetTargetType() TargetType
	Process()
}
type TargetType uint

const (
	JiraType TargetType = iota
)

var (
	targetTypeToString = map[TargetType]string{
		JiraType: "JIRA",
	}
	targetTypeFromString = map[string]TargetType{
		"JIRA": JiraType,
	}
)

func (receiver *TargetType) String() string {
	return targetTypeToString[*receiver]
}

func TargetTypeFromString(s string) TargetType {
	return targetTypeFromString[s]
}

func (receiver *TargetType) MarshalJSON() ([]byte, error) {
	return json.Marshal(targetTypeToString[*receiver])
}

func (receiver *TargetType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*receiver = targetTypeFromString[s]
	return nil
}

type AccessState uint

const (
	NotRegistered AccessState = iota
	Registered
)

type Message struct {
	SourceType SourceType `json:"source_type"`
	TargetType TargetType `json:"target_type"`
	Content    string     `json:"content"`
	UserId     string     `json:"user_id"`
	ChatId     string     `json:"chat_id"`
	MessageId  string     `json:"message_id"`
}

type Comment struct {
	Gid uuid.UUID `json:"gid"`
	*Message
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
	Token      string     `json:"token"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
