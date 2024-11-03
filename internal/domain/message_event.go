package domain

import (
	"github.com/google/uuid"
	"time"
)

type MessageEvent struct {
	Hash       string     `json:"hash"`
	SourceType SourceType `json:"source_type"`
	Text       string     `json:"text"`
	AuthorId   string     `json:"author_id"`
	ChatId     string     `json:"chat_id"`
	MessageId  string     `json:"message_id"`
	Tags       []string   `json:"tags"`
	CommentId  string     `json:"comment_id"`
	CreateTs   time.Time  `json:"create_ts"`
	UpdateTs   time.Time  `json:"update_ts"`
}

func Hash(sourceType SourceType, chatId string, messageId string) string {
	return uuid.NewSHA1(uuid.NameSpaceURL, []byte(sourceType.String()+"_"+chatId+"_"+messageId)).String()
}
