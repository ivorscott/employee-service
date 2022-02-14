// Package msg stores all the events sent or received by the service.
package msg

import "encoding/json"

// UnmarshalMessage parses the JSON-encoded data and returns Message.
func UnmarshalMessage(data []byte) (Message, error) {
	var m Message
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes Message.
func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// UnmarshalMetadata parses the JSON-encoded data and returns Metadata.
func UnmarshalMetadata(data []byte) (Metadata, error) {
	var m Metadata
	err := json.Unmarshal(data, &m)
	return m, err
}

// Marshal JSON encodes Metadata.
func (m *Metadata) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

// Metadata represents additional data about the request.
type Metadata struct {
	TraceID string `json:"traceId"`
	UserID  string `json:"userId"`
}

// Message represents a message in being sent or received.
type Message struct {
	Data     []byte      `json:"data"`
	ID       string      `json:"id"`
	Metadata Metadata    `json:"metadata"`
	Type     MessageType `json:"type"`
}

// MessageType is a type of message. It's either a command or an event.
type MessageType string
