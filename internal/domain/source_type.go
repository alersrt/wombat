package domain

import (
	"encoding/json"
)

type SourceType uint

const (
	Telegram SourceType = iota
)

var (
	sourceTypeToString = map[SourceType]string{
		Telegram: "TELEGRAM",
	}
	sourceTypeFromString = map[string]SourceType{
		"TELEGRAM": Telegram,
	}
)

func (receiver *SourceType) String() string {
	return sourceTypeToString[*receiver]
}

func SourceTypeFromString(s string) SourceType {
	return sourceTypeFromString[s]
}

func (receiver *SourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(sourceTypeToString[*receiver])
}

func (receiver *SourceType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*receiver = sourceTypeFromString[s]
	return nil
}
