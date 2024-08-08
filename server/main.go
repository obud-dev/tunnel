package main

import (
	"os"
	"time"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/model"
	"github.com/obud-dev/tunnel/pkg/svc"
	"github.com/obud-dev/tunnel/pkg/transport"
)

func main() {
	listenOn := os.Getenv("ListenOn")
	api := os.Getenv("Api")
	user := os.Getenv("User")
	password := os.Getenv("Password")

	var server svc.Server
	svcCtx := svc.NewServerCtx(config.ServerConfig{
		ListenOn: listenOn,
		Api:      api,
		User:     user,
		Password: password,
	})

	// 插入测试数据tunnel
	data := &model.Tunnel{ID: "ccf7258f-0e41-4e80-a4ea-18ed8195b98e", Name: "test", Uptime: time.Now().Unix(), Token: "1234abc"}
	_, e := svcCtx.TunnelModel.GetTunnelByID(data.ID)
	if e != nil {
		svcCtx.TunnelModel.Update(data)
	}
	// 插入测试数据route
	route := &model.Route{
		ID:       "1234-5678-abcd-ef",
		TunnelID: "ccf7258f-0e41-4e80-a4ea-18ed8195b98e",
		Hostname: "localhost:5429",
		Prefix:   "/example",
		Target:   "http://0.0.0.0:8080",
		Protocol: model.TypeHttp,
	}

	svcCtx.RouteModel.Update(route)

	server = transport.NewTcpServer(svcCtx)
	go func() {
		server.Listen()
		// if err != nil {
		// 	panic(err)
		// }
	}()
	ApiServer(svcCtx)
}
