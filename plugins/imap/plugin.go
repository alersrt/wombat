package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapclient"
	"github.com/emersion/go-message/charset"
	"github.com/wombat/pkg"
	"io"
	"mime"
	"os"
	"time"
)

type Config struct {
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

// Plugin responsible for getting messages from mail.
type Plugin struct {
	cfg      *Config
	messages chan []byte
}

func New(cfg []byte) (pkg.Plugin, error) {
	imapCfg := &Config{}
	if err := json.Unmarshal(cfg, imapCfg); err != nil {
		return nil, err
	}

	return &Plugin{cfg: imapCfg, messages: make(chan []byte)}, nil
}

func (p *Plugin) Close() error {
	close(p.messages)
	return nil
}

func (p *Plugin) Publish() <-chan []byte {
	return p.messages
}

func (p *Plugin) Run(ctx context.Context) error {
	var client *imapclient.Client

	var debugWriter io.Writer
	if p.cfg.Verbose {
		debugWriter = os.Stdout
	}
	options := &imapclient.Options{
		WordDecoder: &mime.WordDecoder{CharsetReader: charset.Reader},
		DebugWriter: debugWriter,
	}
	client, err := imapclient.DialTLS(p.cfg.Url, options)
	defer func() {
		_ = client.Logout().Wait()
	}()
	if err != nil {
		return fmt.Errorf("mail: read: %v", err)
	}

	if err = client.Login(p.cfg.Username, p.cfg.Password).Wait(); err != nil {
		return fmt.Errorf("mail: read: %v", err)
	}

	if _, err = client.Select(p.cfg.Mailbox, &imap.SelectOptions{ReadOnly: false}).Wait(); err != nil {
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
			time.Sleep(p.cfg.IdleTimeout * time.Millisecond)

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
				bytes, err := json.Marshal(&Message{
					Envelope: item.Envelope,
					Text:     string(item.FindBodySection(&imap.FetchItemBodySection{Specifier: imap.PartSpecifierText})),
				})
				if err != nil {
					return err
				}
				p.messages <- bytes
			}
		}
	}
}
