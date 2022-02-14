// Package message stores all the events sent or received by the service.
package message

import "encoding/json"

func UnmarshalMessage(data []byte) (Message, error) {
	var r Message
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Message) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalMetadata(data []byte) (Metadata, error) {
	var r Metadata
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Metadata) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Metadata struct {
	TraceID string `json:"traceId"`
	UserID  string `json:"userId"`
}

type Message struct {
	Data     interface{} `json:"data"`
	ID       string      `json:"id"`
	Metadata Metadata    `json:"metadata"`
	Type     Type        `json:"type"`
}

type Type string
