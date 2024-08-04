package transport

// 使用tcp协议进行通信

import (
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/message"
	"github.com/obud-dev/tunnel/pkg/model"
	"github.com/obud-dev/tunnel/pkg/svc"
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
	c.conn = conn
	// 发送连接消息
	data, _ := c.conf.Encode()
	c.SendMessage(message.Message{
		Type: message.MessageTypeConnect,
		Data: []byte(data),
	})

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
	ctx      *svc.ServerCtx
	routes   []model.Route        // 路由
	tunnels  map[string]*net.Conn // 隧道ID -> 隧道连接
	messages map[string]*net.Conn // 消息ID -> 外部连接
}

func NewTcpServer(ctx *svc.ServerCtx) *TcpServer {
	tunnels := make(map[string]*net.Conn)
	messages := make(map[string]*net.Conn)
	return &TcpServer{ctx: ctx, tunnels: tunnels, messages: messages}
}

func (s *TcpServer) Listen() error {
	ln, err := net.Listen("tcp", s.ctx.Config.ListenOn)
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
	log.Println("Tunnel connection established with the client")
	channel := make(chan []byte)
	go func() {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Tunnel connection error: %v\n", err)
			return
		}
		log.Printf("Received data on tunnel: %s\n", string(buf[:n]))
		channel <- buf[:n]
	}()
	message_response := &message.Message{}
	message_response.Id = utils.GenerateID()
	message_response.Type = message.MessageTypeDisconnect
	message_response.Data = []byte("connect error")
	if channel != nil {
		message_str := <-channel
		message_request, err := message.Unmarshal(message_str)
		if err != nil {
			log.Println("error unmarshal request message:", err)
			message_response.Data = []byte("data parse failed")
		} else {
			switch message_request.Type {
			case message.MessageTypeConnect:
				log.Println("token:", string(message_request.Data))
				token_request, err := config.ParseFromEncoded(string(message_request.Data))
				if err != nil {
					log.Println("error parse connect token:", err)
				} else {
					// 从数据库中查找tunnel_id对应的记录
					tunnel, err := s.ctx.TunnelModel.GetTunnelByID(token_request.TunnelID)
					if err != nil {
						log.Println("tunnel not found")
						message_response.Data = []byte("tunnel not found in server")
					} else {
						tunnel.Status = "online"
						message_response.Type = message.MessageTypeConnect
						tunnel_json, err := json.Marshal(tunnel)
						if err != nil {
							log.Println("tunnel json marshal failed")
						} else {
							message_response.Data = []byte(string(tunnel_json))
							s.tunnels[tunnel.ID] = &conn
						}
					}
				}
				// 写回数据
				message_byte, _ := message_response.Marshal()
				conn.Write(message_byte)
				// 验证失败，关闭连接
				if message_response.Type == message.MessageTypeDisconnect {
					conn.Close()
				}
				// case message.MessageTypeData:
				// 	fmt.Println("data:", string(m.Data))
				// case message.MessageTypeRouteUpdate:
				// 	fmt.Println("route update")
				// fmt.Println("connected")
			case message.MessageTypeDisconnect:
				fmt.Println("disconnected")
			case message.MessageTypeHeartbeat:
				fmt.Println("heartbeat")
			default:
				fmt.Println("unknown message type")
			}
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
