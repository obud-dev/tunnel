package svc

import (
	"net"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/message"
	"github.com/obud-dev/tunnel/pkg/model"
	"gorm.io/gorm"
)

const (
	DefaultHost     = "0.0.0.0"
	DefaultListenOn = ":5429"
	DefaultApi      = ":8000"
)

type Server interface {
	Listen() error
	HandleConnect(m message.Message, conn net.Conn)
}

type ActiveTunnel struct {
	Conn    net.Conn
	Token   string
	Channel chan []byte
}

type ServerCtx struct {
	Config      config.ServerConfig
	TunnelModel model.TunnelModel
	RouteModel  model.RouteModel
	ServerModel model.ServerModel
	Routes      []model.Route            // 路由
	Tunnels     map[string]*ActiveTunnel // 隧道ID -> 隧道连接
	Messages    map[string]*ActiveTunnel // 消息ID -> 外部连接
	Mutex       sync.Mutex
}

func NewServerCtx(config config.ServerConfig) *ServerCtx {

	if config.Host == "" {
		config.Host = DefaultHost
	}
	if config.ListenOn == "" {
		config.ListenOn = DefaultListenOn
	}
	if config.Api == "" {
		config.Api = DefaultApi
	}
	if config.User == "" {
		config.User = "admin"
	}

	if config.Password == "" {
		// panic("password is required")
		config.Password = "123456"
	}

	db, err := gorm.Open(sqlite.Open("tunnel.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Tunnel{})
	db.AutoMigrate(&model.Route{})
	db.AutoMigrate(&model.Server{})
	tunnelModel := model.NewTunnelModel(db)
	routeModel := model.NewRouteModel(db)
	serverModel := model.NewServerModel(db)

	serverModel.Update(&model.Server{
		Host:     config.Host,
		ListenOn: config.ListenOn,
		Api:      config.Api,
		Version:  "v1.0.0",
	})

	routes, err := routeModel.GetRoutes()
	if err != nil {
		panic(err)
	}

	return &ServerCtx{
		Config:      config,
		TunnelModel: tunnelModel,
		RouteModel:  routeModel,
		ServerModel: serverModel,
		Routes:      routes,
		Tunnels:     map[string]*ActiveTunnel{},
		Messages:    map[string]*ActiveTunnel{},
	}
}
