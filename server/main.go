package main

import (
	"os"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/transport"
)

func main() {
	listenOn := os.Getenv("ListenOn")
	api := os.Getenv("Api")

	if listenOn == "" {
		listenOn = DefaultListenOn
	}
	if api == "" {
		api = DefaultApi
	}

	var server Server
	svcCtx := NewServerCtx(config.ServerConfig{})
	server = transport.NewTcpServer(&svcCtx.Confif)
	err := server.Listen()
	if err != nil {
		panic(err)
	}
	ApiServer(svcCtx)
}
