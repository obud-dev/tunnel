package transport

// 使用tcp协议进行通信

import (
	"fmt"
	"net"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/message"
	"github.com/obud-dev/tunnel/pkg/model"
)

// TCPClient is a client for the TCP protocol
type TcpClient struct {
	conn   net.Conn
	conf   *config.ClientConfig
	routes []model.Route
}

// NewTCPClient creates a new TCP client
func NewTcpClient(token string) *TcpClient {
	// todo: token解析转config
	conf := &config.ClientConfig{}
	routes := make([]model.Route, 0)
	return &TcpClient{conf: conf, routes: routes}
}

func (c *TcpClient) Connect() error {
	// 连接服务器
	// todo: 发送token到服务器
	conn, err := net.Dial("tcp", c.conf.Server)
	if err != nil {
		return err
	}
	c.conn = conn

	for {
		// 读取数据
		buf := make([]byte, 1024)
		n, err := c.conn.Read(buf)
		if err != nil {
			return err
		}
		m, err := message.Unmarshal(buf[:n])
		if err != nil {
			return err
		}
		switch m.Type {
		case message.MessageTypeData:
			fmt.Println("data:", string(m.Data))
			c.RecieveData(*m)
		case message.MessageTypeRouteUpdate:
			fmt.Println("route update")
			c.UpdateRoutes()
		case message.MessageTypeConnect:
			fmt.Println("connected")
		case message.MessageTypeDisconnect:
			fmt.Println("disconnected")
		case message.MessageTypeHeartbeat:
			fmt.Println("heartbeat")
		default:
			fmt.Println("unknown message type")
		}
	}
}

func (c *TcpClient) SendMessage(m message.Message) error {
	data, err := m.Marshal()
	if err != nil {
		return err
	}
	_, err = c.conn.Write(data)
	return err
}

func (c *TcpClient) Close() error {
	return c.conn.Close()
}

func (c *TcpClient) UpdateRoutes() error {
	//todo: 发送消息到服务器，获取路由信息
	return nil
}

func (c *TcpClient) RecieveData(m message.Message) error {
	// todo: 转发数据到内网服务
	// 获取route 通过route知道要转发什么类型到内网host

	// 处理之后的数据返回给服务器
	data := []byte(`处理之后的数据`)
	c.SendMessage(message.Message{
		Type: message.MessageTypeData,
		Data: data,
		Id:   m.Id,
	})

	return nil
}
