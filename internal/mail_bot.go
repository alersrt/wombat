package internal

import (
	"context"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"io"
	"mime"
	"os"
	"time"
)

// MailBot responsible for getting messages from mail.
type MailBot struct {
	cfg *Imap
	out chan *Message
}

func NewMailBot(cfg *Imap) *MailBot {
	return &MailBot{cfg: cfg, out: make(chan *Message)}
}

func (m *MailBot) Close() error {
	close(m.out)
	return nil
}

// Out returns output channel for read
func (m *MailBot) Out() <-chan *Message {
	return m.out
}

// Read starts reading messages from mail.
func (m *MailBot) Read(ctx context.Context) error {
	var client *imapclient.Client

	var debugWriter io.Writer
	if m.cfg.Verbose {
		debugWriter = os.Stdout
	}
	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
		DebugWriter: debugWriter,
	}
	client, err := imapclient.DialTLS(m.cfg.Url, options)
	defer func() {
		_ = client.Logout().Wait()
	}()
	if err != nil {
		return err
	}

	if err = client.Login(m.cfg.Username, m.cfg.Password).Wait(); err != nil {
		return err
	}

	if _, err = client.Select(m.cfg.Mailbox, &imap.SelectOptions{ReadOnly: false}).Wait(); err != nil {
		return err
	}
	defer func() {
		_ = client.Unselect().Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(m.cfg.IdleTimeout * time.Millisecond)

			search, err := client.Search(&imap.SearchCriteria{NotFlag: []imap.Flag{imap.FlagSeen}}, nil).Wait()
			if err != nil {
				return err
			}
			if len(search.AllSeqNums()) == 0 {
				continue
			}
			found, err := client.Fetch(search.All, &imap.FetchOptions{Envelope: true, BodySection: []*imap.FetchItemBodySection{{Specifier: imap.PartSpecifierText}}}).Collect()
			if err != nil {
				return err
			}
			for _, item := range found {
				m.out <- &Message{
					Envelope: item.Envelope,
					Text:     item.FindBodySection(&imap.FetchItemBodySection{Specifier: imap.PartSpecifierText}),
				}
			}
		}
	}
}
