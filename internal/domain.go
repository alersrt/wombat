package internal

import (
	"github.com/emersion/go-imap/v2"
)

type Message struct {
	Text     string         `json:"Text"`
	Envelope *imap.Envelope `json:"Envelope"`
}
