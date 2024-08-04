package main

import (
	"os"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/svc"
	"github.com/obud-dev/tunnel/pkg/transport"
)

func main() {
	listenOn := os.Getenv("ListenOn")
	api := os.Getenv("Api")

	var server svc.Server
	svcCtx := svc.NewServerCtx(config.ServerConfig{
		ListenOn: listenOn,
		Api:      api,
	})
	server = transport.NewTcpServer(svcCtx)
	err := server.Listen()
	if err != nil {
		panic(err)
	}
	ApiServer(svcCtx)
}
