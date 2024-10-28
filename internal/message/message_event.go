package message

type MessageEvent struct {
	SourceType SourceTypeEnum `json:"source_type"`
	Text       string         `json:"text"`
	AuthorId   string         `json:"author_id"`
	ChatId     string         `json:"chat_id"`
	MessageId  string         `json:"message_id"`
}
