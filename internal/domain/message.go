package domain

type Message struct {
	SourceType SourceType `json:"source_type"`
	TargetType TargetType `json:"target_type"`
	Content    string     `json:"content"`
	UserId     string     `json:"user_id"`
	ChatId     string     `json:"chat_id"`
	MessageId  string     `json:"message_id"`
}
