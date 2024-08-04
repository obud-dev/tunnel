package main

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
	UpdateRoutes() error
	HandleConnect(conn net.Conn, m message.Message) error
	HandlePublicData(m message.Message) error
	HandleData(m message.Message) error
	SendMessage(m message.Message) error
}

type ServerCtx struct {
	Confif   config.ServerConfig
	Db       *gorm.DB
	Routes   []model.Route        // 路由
	Tunnels  map[string]*net.Conn // 隧道ID -> 隧道连接
	Messages map[string]*net.Conn // 消息ID -> 外部连接
}

func NewServerCtx(config config.ServerConfig) *ServerCtx {
	db, err := gorm.Open(sqlite.Open("tunnel.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Tunnel{})
	return &ServerCtx{
		Confif:   config,
		Db:       db,
		Routes:   []model.Route{},
		Tunnels:  map[string]*net.Conn{},
		Messages: map[string]*net.Conn{},
	}
}
