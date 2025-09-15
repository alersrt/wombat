package main

import (
	"encoding/json"
	"slices"
	"testing"
	"time"
)

var testedUnit = &Plugin{}

func TestFilter_bool(t *testing.T) {
	cfg := &Config{
		Expr: `self.Text.matches(".*some.*")
&& self.Envelope.exists(f, f == 'Test' || f == 'Check')
`,
	}
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if err := testedUnit.Init(cfgBytes); err != nil {
		t.Fatalf("%+v", err)
	}

	obj, _ := json.Marshal(struct {
		Text     string
		Envelope []string
	}{
		Text:     "some text",
		Envelope: []string{"Test", "Check"},
	})
	actual, err := testedUnit.Eval(obj)
	if err != nil {
		t.Errorf("%+v", err)
	}

	if bT, ok := actual.(bool); !ok || !bT {
		t.Errorf("expected true")
	}
}

func TestFilter_string(t *testing.T) {
	cfg := &Config{
		Expr: `'some' + string(self.Text.matches(".*some.*") && self.Envelope.exists(f, f == 'Test' || f == 'Check'))
`,
	}
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if err := testedUnit.Init(cfgBytes); err != nil {
		t.Fatalf("%+v", err)
	}

	obj, _ := json.Marshal(struct {
		Text     string
		Envelope []string
	}{
		Text:     "some text",
		Envelope: []string{"Test", "Check"},
	})
	actual, err := testedUnit.Eval(obj)
	if err != nil {
		t.Errorf("%+v", err)
	}

	if sT, ok := actual.(string); !ok || sT != "sometrue" {
		t.Errorf("expected [sometrue]")
	}
}

func TestFilter_obj(t *testing.T) {
	cfg := &Config{
		Expr: `{
    "text": self.Text,
    "isCheck": self.Envelope.exists(f, f == 'Check'),
    "envelope": self.Envelope.map(s, {"value": s}),
    "nested": {
        "one": uuid()
    },
    "createdTs": now()
}`,
	}
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if err := testedUnit.Init(cfgBytes); err != nil {
		t.Fatalf("%+v", err)
	}

	obj, _ := json.Marshal(struct {
		Text     string
		Envelope []string
	}{
		Text:     "some text",
		Envelope: []string{"Test", "Check"},
	})
	actual, err := testedUnit.Eval(obj)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	oT, ok := actual.([]byte)
	if !ok {
		t.Fatalf("expected bytes")
	}

	type envelope struct {
		Value string `json:"value"`
	}
	built := &struct {
		Text     string     `json:"text"`
		IsCheck  bool       `json:"isCheck"`
		Envelope []envelope `json:"envelope"`
		Nested   struct {
			One string `json:"one"`
		} `json:"nested"`
		CreatedTs time.Time
	}{}
	if err := json.Unmarshal(oT, built); err != nil {
		t.Fatalf("%+v", err)
	}

	if built.Text != "some text" ||
		!built.IsCheck ||
		len(built.Envelope) != 2 ||
		!slices.Contains(built.Envelope, envelope{Value: "Check"}) ||
		built.Nested.One == "" ||
		built.CreatedTs.Equal(time.Time{}) {
		t.Fatalf("wrong values: %+v", built)
	}
}

func TestFilter_uuid(t *testing.T) {
	cfg := &Config{
		Expr: "uuid(b'00000000-0000-0000-0000-000000000000') == \"00000000-0000-0000-0000-000000000000\"",
	}
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if err := testedUnit.Init(cfgBytes); err != nil {
		t.Fatalf("%+v", err)
	}

	res, err := testedUnit.Eval(nil)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	if actual, ok := res.(bool); !ok || !actual {
		t.Fatalf("not equals")
	}
}
