package message

import "encoding/json"

type SourceType uint

const (
	TELEGRAM SourceType = iota
)

func (receiver SourceType) String() string {
	return [...]string{"TELEGRAM"}[receiver]
}

func (receiver *SourceType) Value(sourceType string) SourceType {
	return map[string]SourceType{"TELEGRAM": TELEGRAM}[sourceType]
}

func (receiver SourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(receiver.String())
}

func (receiver *SourceType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*receiver = receiver.Value(s)
	return nil
}
