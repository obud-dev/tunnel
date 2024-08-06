package transport

// 使用tcp协议进行通信

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

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
		// 读取数据 buf 设置太小会d
		buf := make([]byte, 2048)
		n, err := c.conn.Read(buf)
		if err != nil {
			return err
		}
		fmt.Println("read data:", string(buf[:n]))
		m, err := message.Unmarshal(buf[:n])
		if err != nil {
			return err
		}
		switch m.Type {
		case message.MessageTypeData:
			fmt.Println("data:", string(m.Data))
			go c.RecieveData(*m)
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
	// 转发数据到内网服务
	fmt.Println("client recieve data:", string(m.Data))
	switch m.Protocol {
	case model.TypeHttp:
		// 转发http请求
		reader := bufio.NewReader(bytes.NewReader(m.Data))
		req, err := http.ReadRequest(reader)
		if err != nil {
			fmt.Println("read request error:", err)
			return err
		}
		// todo: 客户端获取的routes
		routes := []model.Route{}
		target := "0.0.0:8080"
		// 通过m.RouteID获取Target,然后传输到Target
		for _, route := range routes {
			if route.ID == m.RouteID {
				target = route.Target
				break
			}
		}
		conn, err := net.Dial("tcp", target)
		if err != nil {
			return err
		}
		// 发送http请求
		req.Write(conn)

		for {
			buf := make([]byte, 2048)
			n, err := conn.Read(buf)
			if err != nil {
				return err
			}
			// 通过tunnel发送本地响应到服务器
			c.SendMessage(message.Message{
				Id:   m.Id,
				Data: buf[:n],
				Type: message.MessageTypeData,
			})
		}
	}
	return nil
}

type TcpServer struct {
	ctx *svc.ServerCtx
}

func NewTcpServer(ctx *svc.ServerCtx) *TcpServer {
	return &TcpServer{ctx: ctx}
}

func (s *TcpServer) Listen() error {
	ln, err := net.Listen("tcp", s.ctx.Config.ListenOn)
	if err != nil {
		return err
	}
	fmt.Printf("Listening on %s\n", s.ctx.Config.ListenOn)
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
	reader := bufio.NewReader(conn)
	for {
		buf := make([]byte, 2048)
		n, err := reader.Read(buf)
		messageStr := buf[:n]
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalln("Tunnel connect error: ", err)
		}
		log.Printf("Received data on tunnel: %s\n", string(messageStr))

		response := &message.Message{}
		response.Id = utils.GenerateID()
		response.Type = message.MessageTypeDisconnect
		response.Data = []byte("connect error")

		request, errx := message.Unmarshal(messageStr)
		if errx != nil {
			// 从外部接收到的数据，转发到内部
			messageId := utils.GenerateID()
			s.ctx.Messages[messageId] = conn
			err := s.HandlePublicData(message.Message{
				Type:     message.MessageTypeData,
				Data:     messageStr,
				Id:       messageId,
				Protocol: model.TypeHttp,
			})
			if err != nil {
				fmt.Println("handle public data error:", err)
				s.ctx.Messages[messageId].Close()
				delete(s.ctx.Messages, messageId)
			}
		} else {
			switch request.Type {
			case message.MessageTypeConnect:
				log.Println("token:", string(request.Data))
				token_request, err := config.ParseFromEncoded(string(request.Data))
				if err != nil {
					log.Println("error parse connect token:", err)
				} else {
					// 从数据库中查找tunnel_id对应的记录
					tunnel, err := s.ctx.TunnelModel.GetTunnelByID(token_request.TunnelID)
					if err != nil {
						log.Println("tunnel not found")
						response.Data = []byte("tunnel not found on server")
					} else {
						tunnel.Status = "online"
						response.Type = message.MessageTypeConnect
						tunnel_json, err := json.Marshal(tunnel)
						if err != nil {
							log.Println("tunnel json marshal failed")
						} else {
							response.Data = []byte(string(tunnel_json))
							s.ctx.Tunnels[tunnel.ID] = conn
						}
					}
				}
				// 写回数据
				message_byte, _ := response.Marshal()
				conn.Write(message_byte)
				// 验证失败，关闭连接
				if response.Type == message.MessageTypeDisconnect {
					conn.Close()
				}
				fmt.Println("connected")
			case message.MessageTypeData:
				if _, ok := s.ctx.Messages[request.Id]; ok {
					s.ctx.Messages[request.Id].Write(request.Data)
				}
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

func (s *TcpServer) HandleConnect(conn net.Conn, m message.Message) error {
	conf, err := config.ParseFromEncoded(string(m.Data))
	if err != nil {
		return err
	}
	// todo: 校验token
	s.ctx.Tunnels[conf.TunnelID] = conn
	return nil
}

func (s *TcpServer) SendMessage(m message.Message) error {
	return nil
}

func (s *TcpServer) HandlePublicData(m message.Message) error {
	tunnelID := "ccf7258f-0e41-4e80-a4ea-18ed8195b98e"
	// todo: 消息规则知道转发到哪个隧道
	// 通过隧道ID获取隧道连接
	if utils.HttpPattern.Match(m.Data) {
		// todo: 处理消息 把消息host 转换为规则host
		// host := GetHostFromHttpMessage(m.Data)
		host := "localhost"
		for _, route := range s.ctx.Routes {
			if route.Hostname == host {
				m.Protocol = route.Protocol
				tunnelID = route.TunnelID
				m.RouteID = route.ID
				break
			}
		}

	}
	if utils.SshPattern.Match(m.Data) {
		m.Protocol = model.TypeSsh
	}

	if _, ok := s.ctx.Tunnels[tunnelID]; !ok {
		return fmt.Errorf("tunnel not found")
	}

	conn := s.ctx.Tunnels[tunnelID]

	mData, err := m.Marshal()
	if err != nil {
		fmt.Println("marshal message error:", err)
		return err
	}
	conn.Write(mData)
	fmt.Println("send data to tunnel")
	return nil
}

func (s *TcpServer) HandleData(m message.Message) error {
	// 通过消息ID获取消息对应

	return nil
}
