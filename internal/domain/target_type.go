package domain

import (
	"encoding/json"
)

type TargetType uint

const (
	JIRA TargetType = iota
)

var (
	targetTypeToString = map[TargetType]string{
		JIRA: "JIRA",
	}
	targetTypeFromString = map[string]TargetType{
		"JIRA": JIRA,
	}
)

func (receiver *TargetType) String() string {
	return targetTypeToString[*receiver]
}

func TargetTypeFromString(s string) TargetType {
	return targetTypeFromString[s]
}

func (receiver *TargetType) MarshalJSON() ([]byte, error) {
	return json.Marshal(targetTypeToString[*receiver])
}

func (receiver *TargetType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*receiver = targetTypeFromString[s]
	return nil
}
