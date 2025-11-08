package main

import "github.com/emersion/go-imap/v2"

func envelopeToEnvelope(value *imap.Envelope) Envelope {
	return Envelope{
		MessageID: value.MessageID,
		InReplyTo: value.InReplyTo,
		Date:      value.Date,
		Subject:   value.Subject,
		From:      addressToAddress(value.From),
		Sender:    addressToAddress(value.Sender),
		ReplyTo:   addressToAddress(value.ReplyTo),
		To:        addressToAddress(value.To),
		Cc:        addressToAddress(value.Cc),
		Bcc:       addressToAddress(value.Bcc),
	}
}

func addressToAddress(value []imap.Address) []Address {
	var addresses []Address
	for _, v := range value {
		addresses = append(addresses, Address{
			Name:    v.Name,
			Mailbox: v.Mailbox,
			Host:    v.Host,
		})
	}
	return addresses
}
