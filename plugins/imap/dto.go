package main

import "time"

// Envelope is the envelope structure of a message.
//
// The subject and addresses are UTF-8 (ie, not in their encoded form). The
// In-Reply-To and Message-ID values contain message identifiers without angle
// brackets.
type Envelope struct {
	Date      time.Time `json:"date"`
	Subject   string    `json:"subject"`
	From      []Address `json:"from"`
	Sender    []Address `json:"sender"`
	ReplyTo   []Address `json:"reply_to"`
	To        []Address `json:"to"`
	Cc        []Address `json:"cc"`
	Bcc       []Address `json:"bcc"`
	InReplyTo []string  `json:"in_reply_to"`
	MessageID string    `json:"message_id"`
}

// Address represents a sender or recipient of a message.
type Address struct {
	Name    string `json:"name"`
	Mailbox string `json:"mailbox"`
	Host    string `json:"host"`
}

// Message represens an output message.
type Message struct {
	Text     string   `json:"text"`
	Envelope Envelope `json:"envelope"`
}
