package domain

import "time"

type MessageEvent struct {
	Hash       string     `json:"hash"`
	EventType  EventType  `json:"event_type"`
	SourceType SourceType `json:"source_type"`
	Text       string     `json:"text"`
	AuthorId   string     `json:"author_id"`
	ChatId     string     `json:"chat_id"`
	MessageId  string     `json:"message_id"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}
