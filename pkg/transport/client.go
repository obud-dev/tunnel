package transport

import "github.com/obud-dev/tunnel/pkg/message"

type TransportClient interface {
	Connect() error
	Close() error
	SendMessage(m message.Message) error
	RecieveData(m message.Message) error
}
