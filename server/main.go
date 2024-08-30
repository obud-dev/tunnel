package main

import (
	"os"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/svc"
	"github.com/obud-dev/tunnel/pkg/transport"
	"github.com/obud-dev/tunnel/pkg/utils"
)

func main() {
	host := os.Getenv("Host")
	listenOn := os.Getenv("ListenOn")
	api := os.Getenv("Api")
	user := os.Getenv("User")
	password := os.Getenv("Password")

	utils.InitLogger()

	var server svc.Server
	svcCtx := svc.NewServerCtx(config.ServerConfig{
		Host:     host,
		ListenOn: listenOn,
		Api:      api,
		User:     user,
		Password: password,
	})

	// go utils.PrintMemoryUsage()

	server = transport.NewTcpServer(svcCtx)

	go server.Listen()
	go ApiServer(svcCtx)

	select {}
}
