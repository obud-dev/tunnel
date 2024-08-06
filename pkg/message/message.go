package message

import (
	"encoding/json"

	"github.com/obud-dev/tunnel/pkg/model"
)

type MessageType int

const (
	MessageTypeData MessageType = iota
	MessageTypeConnect
	MessageTypeDisconnect
	MessageTypeHeartbeat
)

type Message struct {
	Type     MessageType    `json:"type"`
	Data     []byte         `json:"data"`
	Id       string         `json:"id"`
	Protocol model.Protocol `json:"protocol"`
	RouteID  string         `json:"routeId"`
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func Unmarshal(data []byte) (*Message, error) {
	m := &Message{}
	err := json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
