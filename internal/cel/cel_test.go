package cel

import (
	"encoding/json"
	"slices"
	"testing"
	"time"
)

func TestFilter_bool(t *testing.T) {
	expr := `self.Text.matches(".*some.*")
&& self.Envelope.exists(f, f == 'Test' || f == 'Check')
`

	testedUnit, err := NewCel(expr)
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

	actual, err := testedUnit.EvalBytes(obj)
	if err != nil {
		t.Errorf("%+v", err)
	}

	var bT bool
	err = json.Unmarshal(actual, &bT)
	if err != nil || !bT {
		t.Fatalf("not equals: err=%+v, act=%v", err, bT)
	}
}

func TestFilter_string(t *testing.T) {
	expr := `'some' + string(self.Text.matches(".*some.*") && self.Envelope.exists(f, f == 'Test' || f == 'Check'))
`
	testedUnit, err := NewCel(expr)
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
	actual, err := testedUnit.EvalBytes(obj)
	if err != nil {
		t.Errorf("%+v", err)
	}

	var sT string
	err = json.Unmarshal(actual, &sT)
	if err != nil || sT != "sometrue" {
		t.Errorf("expected [sometrue]: act=%v", sT)
	}
}

func TestFilter_obj(t *testing.T) {
	expr := `{
    "text": self.Text,
    "isCheck": self.Envelope.exists(f, f == 'Check'),
    "envelope": self.Envelope.map(s, {"value": s}),
    "nested": {
        "one": uuid()
    },
    "createdTs": now()
}`

	testedUnit, err := NewCel(expr)
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
	actual, err := testedUnit.EvalBytes(obj)
	if err != nil {
		t.Fatalf("%+v", err)
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
	if err := json.Unmarshal(actual, built); err != nil {
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
	expr := "uuid(b'00000000-0000-0000-0000-000000000000') == \"00000000-0000-0000-0000-000000000000\""

	testedUnit, err := NewCel(expr)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	actual, err := testedUnit.EvalBytes(nil)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	var bT bool
	err = json.Unmarshal(actual, &bT)
	if err != nil || !bT {
		t.Fatalf("not equals: err=%+v, act=%v", err, actual)
	}
}
