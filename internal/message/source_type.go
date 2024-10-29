package message

import (
	"encoding/json"
)

type SourceType uint

const (
	TELEGRAM SourceType = iota
)

var (
	toString = map[SourceType]string{
		TELEGRAM: "TELEGRAM",
	}
	fromString = map[string]SourceType{
		"TELEGRAM": TELEGRAM,
	}
)

func (receiver SourceType) String() string {
	return toString[receiver]
}

func (receiver *SourceType) FromString(sourceType string) SourceType {
	return fromString[sourceType]
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
	*receiver = receiver.FromString(s)
	return nil
}
