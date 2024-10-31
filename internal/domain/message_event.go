package domain

type MessageEvent struct {
	EventType  EventType  `json:"event_type"`
	SourceType SourceType `json:"source_type"`
	Text       string     `json:"text"`
	AuthorId   string     `json:"author_id"`
	ChatId     string     `json:"chat_id"`
	MessageId  string     `json:"message_id"`
}
