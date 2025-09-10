package internal

import (
	"context"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"mime"
	"time"
)

type App struct {
}

func (a *App) Init() error {
	return nil
}

const (
	address     = "imap.example.org:993"
	username    = "username"
	password    = "password"
	mailbox     = "INBOX"
	idleTimeout = time.Minute
	queueSize   = 1000
)

func (a *App) Do(ctx context.Context) error {
	buffer := make(chan any)
	defer close(buffer)

	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
		UnilateralDataHandler: &imapclient.UnilateralDataHandler{
			Fetch: func(msg *imapclient.FetchMessageData) {

			},
		},
	}
	client, err := imapclient.DialTLS(address, options)
	if err != nil {
		return err
	}

	if err = client.Login(username, password).Wait(); err != nil {
		return err
	}
	defer client.Logout().Wait()

	if _, err = client.Select(mailbox, nil).Wait(); err != nil {
		return err
	}
	defer client.Unselect().Wait()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			idle, err := client.Idle()
			if err != nil {
				continue
			}
			time.Sleep(idleTimeout)
			if err := idle.Close(); err != nil {
				continue
			}
		}
	}
}
