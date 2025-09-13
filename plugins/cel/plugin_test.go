package main

import (
	"encoding/json"
	"testing"
)

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

	testedUnit, err := New(cfgBytes)
	if err != nil {
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

	testedUnit, err := New(cfgBytes)
	if err != nil {
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
        "one": 1
    }
}`,
	}
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	testedUnit, err := New(cfgBytes)
	if err != nil {
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

	built := &struct {
		Text     string `json:"text"`
		IsCheck  bool   `json:"isCheck"`
		Envelope []struct {
			Value string `json:"value"`
		} `json:"envelope"`
		Nested struct {
			One int `json:"one"`
		} `json:"nested"`
	}{}
	if err := json.Unmarshal(oT, built); err != nil {
		t.Fatalf("%+v", err)
	}

	if built.Text != "some text" || !built.IsCheck || len(built.Envelope) != 2 || built.Nested.One != 1 {
		t.Fatalf("wrong values: %+v", built)
	}
}
