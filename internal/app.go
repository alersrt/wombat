package internal

import (
	"context"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"mime"
)

type App struct {
}

func (a *App) Init() error {

	return nil
}

func (a *App) Do(ctx context.Context) error {
	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
	}
	client, err := imapclient.DialTLS("imap.example.org:993", options)
	if err != nil {
		return err
	}
	defer client.Logout()

	client.Login("username", "password")

	<-ctx.Done()
	return nil
}
