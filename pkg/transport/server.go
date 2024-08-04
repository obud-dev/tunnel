package transport

import (
	"net"

	"github.com/obud-dev/tunnel/pkg/message"
)

type TransportServer interface {
	Listen() error
	UpdateRoutes() error
	HandleConnect(conn net.Conn, m message.Message) error
	HandlePublicData(m message.Message) error
	HandleData(m message.Message) error
	SendMessage(m message.Message) error
}
