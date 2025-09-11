package filter

import (
	"testing"
	"wombat/internal"
)

func TestFilter(t *testing.T) {
	filterContent := `Value.some == 0
&& Headers.key1 in ['value']
&& Headers.size() != 0
&& string(Key).matches(".*test_key.*")
&& Timestamp > timestamp('1970-01-01T00:00:00.000Z')
`
	testedUnit, err := NewFilter(filterContent)
	if err != nil {
		t.Fatalf("%+v", err)
	}

	msg := &internal.Message{}
	ok, err := testedUnit.Eval(msg)
	if err != nil {
		t.Errorf("%+v", err)
	}
	if !ok {
		t.Errorf("expected true")
	}
}

func BenchmarkFilter_Eval(b *testing.B) {
	filterContent := `Value.some == 0
&& Headers.key1 in ['value']
&& Headers.size() != 0
&& Key.kf == 'test_key'
&& Timestamp > timestamp('1970-01-01T00:00:00.000Z')
`
	testedUnit, _ := NewFilter(filterContent)

	msg := &internal.Message{}

	for i := 0; i < b.N; i++ {
		_, _ = testedUnit.Eval(msg)
	}
}
