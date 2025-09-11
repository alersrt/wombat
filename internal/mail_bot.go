package internal

import (
	"context"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"mime"
	"os"
	"time"
)

type MailBot struct {
	cfg *Imap
}

func NewMailBot(cfg *Imap) *MailBot {
	return &MailBot{cfg: cfg}
}

func (a *MailBot) Read(ctx context.Context, dest chan any) error {
	var client *imapclient.Client

	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
		DebugWriter: os.Stdout,
	}
	client, err := imapclient.DialTLS(a.cfg.Address, options)
	defer func() {
		_ = client.Logout().Wait()
	}()
	if err != nil {
		return err
	}

	if err = client.Login(a.cfg.Username, a.cfg.Password).Wait(); err != nil {
		return err
	}

	if _, err = client.Select(a.cfg.Mailbox, &imap.SelectOptions{ReadOnly: false}).Wait(); err != nil {
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
			time.Sleep(a.cfg.IdleTimeout * time.Millisecond)

			search, err := client.Search(&imap.SearchCriteria{NotFlag: []imap.Flag{imap.FlagSeen}}, nil).Wait()
			if err != nil {
				continue
			}
			if len(search.AllSeqNums()) == 0 {
				continue
			}
			found, err := client.Fetch(search.All, &imap.FetchOptions{Envelope: true, BodySection: []*imap.FetchItemBodySection{{Specifier: imap.PartSpecifierText}}}).Collect()
			if err != nil {
				continue
			}
			for _, item := range found {
				dest <- item
				//fmt.Printf("==============================================================\n")
				//fmt.Printf("%s\n", string(item.FindBodySection(&imap.FetchItemBodySection{Specifier: imap.PartSpecifierText})))
				//fmt.Printf("%+v\n", item.Envelope)
				//fmt.Printf("==============================================================\n")
			}
		}
	}
}
