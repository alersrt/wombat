package internal

import "github.com/emersion/go-imap/v2"

type Message struct {
	Text     []byte         `json:"text"`
	Envelope *imap.Envelope `json:"envelope"`
}
