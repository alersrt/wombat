package domain

import (
	"encoding/json"
)

type SourceType uint

const (
	TELEGRAM SourceType = iota
)

var (
	sourceTypeToString = map[SourceType]string{
		TELEGRAM: "TELEGRAM",
	}
	sourceTypeFromString = map[string]SourceType{
		"TELEGRAM": TELEGRAM,
	}
)

func (receiver *SourceType) String() string {
	return sourceTypeToString[*receiver]
}

func (receiver *SourceType) FromString(s string) {
	*receiver = sourceTypeFromString[s]
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
