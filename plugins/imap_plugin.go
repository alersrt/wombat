package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"io"
	"mime"
	"os"
	"time"
)

type ImapCfg struct {
	Url         string        `yaml:"url"`
	Username    string        `yaml:"username"`
	Password    string        `yaml:"password"`
	Mailbox     string        `yaml:"mailbox"`
	IdleTimeout time.Duration `yaml:"idleTimeout"`
	Verbose     bool          `yaml:"verbose"`
}

type Message struct {
	Text     string         `json:"Text"`
	Envelope *imap.Envelope `json:"Envelope"`
}

// ImapSrc responsible for getting messages from mail.
type ImapSrc struct {
	cfg      *ImapCfg
	messages chan any
}

func New(cfg map[string]any) (*ImapSrc, error) {
	bytes, err := json.Marshal(cfg)
	if err != nil {
		return nil, err
	}
	imapCfg := &ImapCfg{}
	if err := json.Unmarshal(bytes, imapCfg); err != nil {
		return nil, err
	}

	return &ImapSrc{cfg: imapCfg, messages: make(chan any)}, nil
}

func (s *ImapSrc) Close() error {
	close(s.messages)
	return nil
}

// Publish returns output channel for read
func (s *ImapSrc) Publish() <-chan any {
	return s.messages
}

// Run starts reading messages from mail.
func (s *ImapSrc) Run(ctx context.Context) error {
	var client *imapclient.Client

	var debugWriter io.Writer
	if s.cfg.Verbose {
		debugWriter = os.Stdout
	}
	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
		DebugWriter: debugWriter,
	}
	client, err := imapclient.DialTLS(s.cfg.Url, options)
	defer func() {
		_ = client.Logout().Wait()
	}()
	if err != nil {
		return fmt.Errorf("mail: read: %v", err)
	}

	if err = client.Login(s.cfg.Username, s.cfg.Password).Wait(); err != nil {
		return fmt.Errorf("mail: read: %v", err)
	}

	if _, err = client.Select(s.cfg.Mailbox, &imap.SelectOptions{ReadOnly: false}).Wait(); err != nil {
		return fmt.Errorf("mail: read: %v", err)
	}
	defer func() {
		_ = client.Unselect().Wait()
	}()

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			time.Sleep(s.cfg.IdleTimeout * time.Millisecond)

			search, err := client.Search(&imap.SearchCriteria{NotFlag: []imap.Flag{imap.FlagSeen}}, nil).Wait()
			if err != nil {
				return fmt.Errorf("mail: read: %v", err)
			}
			if len(search.AllSeqNums()) == 0 {
				continue
			}
			found, err := client.Fetch(search.All, &imap.FetchOptions{Envelope: true, BodySection: []*imap.FetchItemBodySection{{Specifier: imap.PartSpecifierText}}}).Collect()
			if err != nil {
				return fmt.Errorf("mail: read: %v", err)
			}
			for _, item := range found {
				s.messages <- &Message{
					Envelope: item.Envelope,
					Text:     string(item.FindBodySection(&imap.FetchItemBodySection{Specifier: imap.PartSpecifierText})),
				}
			}
		}
	}
}
