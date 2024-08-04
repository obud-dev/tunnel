package transport

// 使用tcp协议进行通信

import (
	"fmt"
	"net"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/message"
	"github.com/obud-dev/tunnel/pkg/model"
	"github.com/obud-dev/tunnel/pkg/utils"
)

// TCPClient is a client for the TCP protocol
type TcpClient struct {
	conn net.Conn
	conf *config.ClientConfig
}

// NewTCPClient creates a new TCP client
func NewTcpClient(token string) *TcpClient {
	conf, err := config.ParseFromEncoded(token)
	if err != nil {
		panic(err)
	}
	return &TcpClient{conf: conf}
}

func (c *TcpClient) Connect() error {
	// 连接服务器
	conn, err := net.Dial("tcp", c.conf.Server)
	if err != nil {
		return err
	}
	// 发送连接消息
	data, _ := c.conf.Encode()
	c.SendMessage(message.Message{
		Type: message.MessageTypeConnect,
		Data: []byte(data),
	})

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
		// case message.MessageTypeRouteUpdate:
		// 	fmt.Println("route update")
		// 	c.UpdateRoutes()
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

type TcpServer struct {
	conf     *config.ServerConfig // 服务器配置
	routes   []model.Route        // 路由
	tunnels  map[string]*net.Conn // 隧道ID -> 隧道连接
	messages map[string]*net.Conn // 消息ID -> 外部连接
}

func NewTcpServer(conf *config.ServerConfig) *TcpServer {
	return &TcpServer{conf: conf}
}

func (s *TcpServer) Listen() error {
	ln, err := net.Listen("tcp", s.conf.ListenOn)
	if err != nil {
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *TcpServer) handleConn(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		m, err := message.Unmarshal(buf[:n])
		if err != nil {
			// 从外部接收到的数据，转发到内部
			messageId := utils.GenerateID()
			s.messages[messageId] = &conn
			s.SendMessage(message.Message{
				Type: message.MessageTypeData,
				Data: buf[:n],
				Id:   messageId,
			})
			return
		}
		// 处理消息
		switch m.Type {
		case message.MessageTypeData:
			fmt.Println("data:", string(m.Data))
		// case message.MessageTypeRouteUpdate:
		// 	fmt.Println("route update")
		case message.MessageTypeConnect:
			// todo: 处理连接消息， 校验token
			// 安装，如果隧道ID已经存在，关闭之前的连接
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

func (s *TcpServer) UpdateRoutes() error {
	// todo: 更新路由
	s.routes = []model.Route{}
	return nil
}

func (s *TcpServer) HandleConnect(conn net.Conn, m message.Message) error {
	conf, err := config.ParseFromEncoded(string(m.Data))
	if err != nil {
		return err
	}
	// todo: 校验token
	s.tunnels[conf.TunnelID] = &conn
	return nil
}

func (s *TcpServer) SendMessage(m message.Message) error {
	return nil
}

func (s *TcpServer) HandlePublicData(m message.Message) error {
	// todo: 消息规则知道转发到哪个隧道
	// 通过隧道ID获取隧道连接
	// tunnelID := ""
	// conn := *s.tunnels[tunnelID]
	// // todo: 处理消息 把消息host 转换为规则host
	// data := []byte(`处理之后的数据`)
	// m.Data = data
	// // 发送消息
	// mData, err := m.Marshal()
	// if err != nil {
	// 	return err
	// }
	// conn.Write(mData)
	// return err
	return nil
}

func (s *TcpServer) HandleData(m message.Message) error {
	// 通过消息ID获取消息对应

	return nil
}
