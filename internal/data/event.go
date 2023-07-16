package data

import (
	"encoding/json"
	"time"
)

const DateTimeLayout = "2006-01-02 15:04:05.000000"

// A type containg time.Time that implements the Unmarshaler interface
type UnmarshalerTime struct {
	time.Time
}

func (u *UnmarshalerTime) UnmarshalJSON(b []byte) error {
	var timestamp string
	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return err
	}
	u.Time, err = time.Parse(DateTimeLayout, timestamp)
	return err
}

//Data: Here the input and output data structures are stored
// along with their associated methods

type Event struct {
	Timestamp       UnmarshalerTime
	Translation_id  string
	Source_language string
	Target_language string
	Client_name     string
	Event_name      string
	Nr_words        uint
	Duration        uint
}

func EventFromString(jsonStr string) (*Event, error) {

	var event Event

	err := json.Unmarshal([]byte(jsonStr), &event)

	if err != nil {
		return nil, err
	}
	return &event, nil
}
