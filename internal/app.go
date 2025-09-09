package internal

import (
	"context"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"mime"
)

type App struct {
}

func (a *App) Init() error {
	return nil
}

const (
	address   = "imap.example.org:993"
	username  = "username"
	password  = "password"
	mailbox   = "INBOX"
	queueSize = 1000
)

func (a *App) Do(ctx context.Context) error {

	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
	}
	client, err := imapclient.DialTLS(address, options)
	if err != nil {
		return err
	}
	defer client.Logout().Wait()

	err = client.Login(username, password).Wait()
	if err != nil {
		return err
	}

	_, err = client.Select(mailbox, nil).Wait()
	if err != nil {
		return err
	}
	defer client.Unselect().Wait()

	criteria := &imap.SearchCriteria{
		NotFlag: []imap.Flag{imap.FlagSeen},
	}

	buffer := make(chan any)
	defer close(buffer)

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			search, err := client.Search(criteria, nil).Wait()
			if err != nil {
				continue
			}

			if client.Mailbox().NumMessages > 0 {
				seqSet := new(imap.SeqSet)
				seqSet.AddRange(search.Min, search.Max)
				fetchOptions := &imap.FetchOptions{Envelope: true}
				messages, err := client.Fetch(seqSet, fetchOptions).Collect()
				if err != nil {
					continue
				}

				for _, m := range messages {
					buffer <- m
				}
			}
		}
	}
}
