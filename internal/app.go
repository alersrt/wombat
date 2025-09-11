package internal

import (
	"context"
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"log/slog"
	"mime"
	"os"
	"time"
)

type App struct {
	cfg *Config
}

func (a *App) Init(cfg *Config) error {
	a.cfg = cfg
	return nil
}

func (a *App) Do(ctx context.Context) error {
	var client *imapclient.Client

	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
		DebugWriter: os.Stdout,
	}
	client, err := imapclient.DialTLS(a.cfg.Imap.Address, options)
	defer func() {
		_ = client.Logout().Wait()
	}()
	if err != nil {
		return err
	}

	if err = client.Login(a.cfg.Imap.Username, a.cfg.Imap.Password).Wait(); err != nil {
		return err
	}

	if _, err = client.Select(a.cfg.Imap.Mailbox, &imap.SelectOptions{ReadOnly: false}).Wait(); err != nil {
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
			time.Sleep(a.cfg.Imap.IdleTimeout * time.Millisecond)

			search, err := client.Search(&imap.SearchCriteria{NotFlag: []imap.Flag{imap.FlagSeen}}, nil).Wait()
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				continue
			}

			found, err := client.Fetch(imap.UIDSetNum(search.AllUIDs()...), nil).Collect()
			if err != nil {
				slog.Error(fmt.Sprintf("%+v", err))
				continue
			}
			for _, item := range found {
				fmt.Printf("==============================================================\n")
				fmt.Printf("%s\n", string(item.FindBodySection(&imap.FetchItemBodySection{Specifier: imap.PartSpecifierText})))
				fmt.Printf("==============================================================\n")
			}
		}
	}
}
