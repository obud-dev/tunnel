package message

import "encoding/json"

type MessageType int

const (
	MessageTypeData MessageType = iota
	MessageTypeConnect
	MessageTypeDisconnect
	MessageTypeHeartbeat
	MessageTypeRouteUpdate
)

type Message struct {
	Type    MessageType `json:"type"`
	Data    []byte      `json:"data"`
	RouteID string      `json:"route_id"`
	Id      string      `json:"id"`
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
