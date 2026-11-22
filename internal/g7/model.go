// Package g7 holds the emulator of a GANTNER G7 device and the models of its messages.
//
// The models are the Go form of CommandDto.cs in the web-services project. They cover the data
// objects ("Data") of the GANTNER GT7 App API, version 1.15.0, and of the G7 Time and Attendance
// App API, version 1.7.0.
//
// Names: "...RequestBody" is the data of a request, "...ResponseBody" the data of a response and
// "...EventBody" the data of an event. A type without one of these endings is a part of other
// data objects or has more than one use. A command that is not here has an empty data object.
//
// An optional field is a pointer, a slice or a map, and has the "omitzero" tag. The JSON does not
// contain the field when it is nil. Use new(value) to set a pointer field, for example
// new("INTERNAL") or new(3000).
package g7

import (
	"encoding/json"
	"time"
)

// Command is the envelope of each message between the host and the device.
type Command struct {
	CommandType CommandType `json:"Cmd"`
	MessageType MessageType `json:"MT"`

	// TransactionID connects a response to its request. An event has no transaction ID.
	TransactionID *int64 `json:"TID,omitzero"`

	GlobalMessageID *string `json:"GID,omitzero"`

	// Data is the data object of the command. Use "{}" for an empty data object.
	Data json.RawMessage `json:"Data"`

	// State is the result of a response. Requests and events usually have no state.
	State *State `json:"State,omitzero"`
}

func NewCommand(cmd CommandType, mt MessageType, body any) (Command, error) {
	data := json.RawMessage("{}")
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return Command{}, err
		}
		data = encoded
	}
	return Command{CommandType: cmd, MessageType: mt, Data: data}, nil
}

func DecodeData[T any](c Command) (T, error) {
	var body T
	err := json.Unmarshal(c.Data, &body)
	return body, err
}

// MessageType is the "MT" value of a message.
type MessageType string

const (
	MessageTypeRequest  MessageType = "Req"
	MessageTypeResponse MessageType = "Rsp"
	MessageTypeEvent    MessageType = "Evt"
)

// State is the "State" value of a response. The JSON value is a number.
type State int

const (
	StateOK             State = 0
	StateGeneralError   State = 1
	StateParameterError State = 2

	// StateResponseTooLarge: the response is larger than the maximum message size (Amazon API
	// Gateway only).
	StateResponseTooLarge State = 253
	StateLoginRequired    State = 254
	StateUnknownCommand   State = 255
)

// HeartbeatBody is the data of the Heartbeat response.
type HeartbeatBody struct {
	// HeartbeatInterval is in seconds.
	HeartbeatInterval int       `json:"HBI"`
	CurrentTime       time.Time `json:"RT"`
}

// DateTime is a time that the host can write with or without a UTC offset, for example
// 2026-01-01T00:00:00 or 2026-01-01T00:00:00Z. A time without an offset reads as UTC.
type DateTime struct {
	time.Time
}

const dateTimeNoOffset = "2006-01-02T15:04:05.999999999"

// UnmarshalJSON reads a time with or without a UTC offset.
func (d *DateTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var text string
	if err := json.Unmarshal(data, &text); err != nil {
		return err
	}

	parsed, err := time.Parse(time.RFC3339Nano, text)
	if err != nil {
		parsed, err = time.Parse(dateTimeNoOffset, text)
	}
	if err != nil {
		return err
	}

	d.Time = parsed
	return nil
}
