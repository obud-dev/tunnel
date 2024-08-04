package svc

import (
	"net"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/message"
	"github.com/obud-dev/tunnel/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	DefaultListenOn = ":5429"
	DefaultApi      = ":8000"
)

type Server interface {
	Listen() error
	HandleConnect(conn net.Conn, m message.Message) error
	HandlePublicData(m message.Message) error
	HandleData(m message.Message) error
	SendMessage(m message.Message) error
}

type ServerCtx struct {
	Config      config.ServerConfig
	TunnelModel model.TunnelModel
	RouteModel  model.RouteModel
	Routes      []model.Route       // 路由
	Tunnels     map[string]net.Conn // 隧道ID -> 隧道连接
	Messages    map[string]net.Conn // 消息ID -> 外部连接
}

func NewServerCtx(config config.ServerConfig) *ServerCtx {

	if config.ListenOn == "" {
		config.ListenOn = DefaultListenOn
	}
	if config.Api == "" {
		config.Api = DefaultApi
	}

	db, err := gorm.Open(sqlite.Open("tunnel.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Tunnel{})
	return &ServerCtx{
		Config:      config,
		TunnelModel: model.NewTunnelModel(db),
		RouteModel:  model.NewRouteModel(db),
		Routes:      []model.Route{},
		Tunnels:     map[string]net.Conn{},
		Messages:    map[string]net.Conn{},
	}
}
