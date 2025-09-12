package internal

import (
	"github.com/emersion/go-imap/v2"
	"testing"
	"time"
	"wombat/plugins"
)

func TestFilter(t *testing.T) {
	filterContent := `self.Text.matches(".*some.*")
&& self.Envelope.From.exists(f, f.Name == 'Test')
`
	testedUnit, err := NewFilter(filterContent)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	msg := &main.Message{
		Text: "some test",
		Envelope: &imap.Envelope{
			Date:    time.Time{},
			Subject: "Some test subj",
			From: []imap.Address{{
				Name:    "Test",
				Mailbox: "usertest",
				Host:    "test.dev",
			}},
			Sender: []imap.Address{{
				Name:    "Test",
				Mailbox: "usertest",
				Host:    "test.dev",
			}},
			ReplyTo:   nil,
			To:        nil,
			Cc:        nil,
			Bcc:       nil,
			InReplyTo: nil,
			MessageID: "",
		},
	}
	ok, err := testedUnit.Eval(msg)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if !ok {
		t.Errorf("expected true")
	}
}

func BenchmarkFilter_Eval(b *testing.B) {
	filterContent := `self.Text.matches(".*some.*")
&& self.Envelope.From.exists(f, f.Name == 'Test')
`
	testedUnit, _ := NewFilter(filterContent)

	msg := &main.Message{
		Text: "some test",
		Envelope: &imap.Envelope{
			Date:    time.Time{},
			Subject: "Some test subj",
			From: []imap.Address{{
				Name:    "Test",
				Mailbox: "usertest",
				Host:    "test.dev",
			}},
			Sender: []imap.Address{{
				Name:    "Test",
				Mailbox: "usertest",
				Host:    "test.dev",
			}},
			ReplyTo:   nil,
			To:        nil,
			Cc:        nil,
			Bcc:       nil,
			InReplyTo: nil,
			MessageID: "",
		},
	}

	for i := 0; i < b.N; i++ {
		_, _ = testedUnit.Eval(msg)
	}
}
