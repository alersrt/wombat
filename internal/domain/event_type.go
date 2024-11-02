package domain

import "encoding/json"

type EventType uint

const (
	CREATE EventType = iota
	UPDATE
)

var (
	eventTypeToString = map[EventType]string{
		CREATE: "CREATE",
		UPDATE: "UPDATE",
	}
	eventTypeFromString = map[string]EventType{
		"CREATE": CREATE,
		"UPDATE": UPDATE,
	}
)

func (receiver *EventType) String() string {
	return eventTypeToString[*receiver]
}

func EventTypeFromString(s string) EventType {
	return eventTypeFromString[s]
}

func (receiver *EventType) MarshalJSON() ([]byte, error) {
	return json.Marshal(eventTypeToString[*receiver])
}

func (receiver *EventType) UnmarshalJSON(b []byte) error {
	var s string
	err := json.Unmarshal(b, &s)
	if err != nil {
		return err
	}
	*receiver = eventTypeFromString[s]
	return nil
}
