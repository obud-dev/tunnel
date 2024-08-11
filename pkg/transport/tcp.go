package transport

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/message"
	"github.com/obud-dev/tunnel/pkg/model"
	"github.com/obud-dev/tunnel/pkg/svc"
	"github.com/obud-dev/tunnel/pkg/utils"
)

// TCP Client Constants
const (
	heartbeatInterval = 10 * time.Second // 心跳包发送间隔
	reconnectAttempts = 3                // 最大重连尝试次数
	reconnectInterval = 10 * time.Second // 每次重连间隔
	heartTimeout      = 5 * time.Second  // 心跳包超时时间
)

// TcpClient is a client for the TCP protocol
type TcpClient struct {
	conn net.Conn
	conf *config.ClientConfig
}

// NewTcpClient creates a new TCP client
func NewTcpClient(token string) (*TcpClient, error) {
	conf, err := config.ParseFromEncoded(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &TcpClient{conf: conf}, nil
}

// Connect establishes a TCP connection to the server
func (c *TcpClient) Connect() error {
	var err error
	c.conn, err = net.Dial("tcp", c.conf.Server)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	log.Info().Msgf("Connecting to server: %s", c.conf.Server)
	// Send connection message
	data, err := c.conf.Encode()
	if err != nil {
		return fmt.Errorf("failed to encode config: %w", err)
	}
	c.SendMessage(message.Message{
		Type: message.MessageTypeConnect,
		Data: []byte(data),
	})

	go c.readLoop()
	go c.Heartbeat()
	select {}
}

// readLoop handles incoming messages from the server
func (c *TcpClient) readLoop() {
	defer c.conn.Close()

	reader := bufio.NewReader(c.conn)
	for {
		buf := make([]byte, 2048)
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Info().Msg("Connection closed by server")
			} else {
				log.Error().Err(err).Msg("Error reading from server")
			}
			break
		}

		m, err := message.Unmarshal(buf[:n])
		if err != nil {
			log.Error().Err(err).Msgf("Failed to unmarshal message %s", buf[:n])
			continue
		}

		c.handleMessage(m)
	}
}

// handleMessage processes different types of messages from the server
func (c *TcpClient) handleMessage(m *message.Message) {
	switch m.Type {
	case message.MessageTypeData:
		go c.RecieveData(*m)
	case message.MessageTypeConnect:
		log.Info().Msg("Connected to server")
	case message.MessageTypeDisconnect:
		log.Info().Msg("Disconnected from server")
	case message.MessageTypeHeartbeat:
		log.Debug().Msg("Received heartbeat")
	default:
		log.Warn().Msg("Unknown message type")
	}
}

// SendMessage sends a message to the server
func (c *TcpClient) SendMessage(m message.Message) {
	data, err := m.Marshal()
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return
	}

	_, err = c.conn.Write(data)
	if err != nil {
		log.Error().Err(err).Msgf("Failed to send message: %v", m.Type)
		c.ReconnectToServer()
		return
	}
}

// RecieveData processes data received from the server
func (c *TcpClient) RecieveData(m message.Message) {
	log.Debug().Msg("Received data from server")

	// Handle the received data based on the protocol
	switch m.Protocol {
	case model.TypeHttp:
		// For HTTP data, initiate a request
		c.handleHttpData(m)
	default:
		log.Warn().Msg("Unknown protocol type for received data")
	}
}

// handleHttpData processes HTTP data received from the server
func (c *TcpClient) handleHttpData(m message.Message) {
	// You need to implement the actual HTTP handling logic here.
	// For example, you might create a new HTTP request using the received data.
	conn, err := net.Dial("tcp", m.Target)
	if err != nil {
		log.Error().Err(err).Msg("Error connecting to target")
		return
	}
	conn.Write(m.Data)
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			log.Error().Err(err).Msg("Error reading from target")
			conn.Close()
			break
		}
		log.Info().Msg("Received response")
		// 通过tunnel发送本地响应到服务器
		c.SendMessage(message.Message{
			Id:   m.Id,
			Data: buf[:n],
			Type: message.MessageTypeData,
		})
	}
}

// Heartbeat sends periodic heartbeat messages to the server
func (c *TcpClient) Heartbeat() {
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for range ticker.C {
		c.conn.SetWriteDeadline(time.Now().Add(heartTimeout))
		m := message.Message{
			Type: message.MessageTypeHeartbeat,
			Data: []byte("ping"),
		}
		mBytes, _ := m.Marshal()
		_, err := c.conn.Write(mBytes)
		if err != nil {
			log.Error().Err(err).Msg("Failed to send heartbeat")
			c.ReconnectToServer()
		}
	}
}

// ReconnectToServer attempts to reconnect to the TCP server
func (c *TcpClient) ReconnectToServer() error {
	for i := 0; i < reconnectAttempts; i++ {
		log.Info().Msgf("Reconnect attempt %d/%d", i+1, reconnectAttempts)
		if err := c.Connect(); err == nil {
			return nil
		}
		time.Sleep(reconnectInterval)
	}
	return fmt.Errorf("all reconnect attempts failed")
}

// TcpServer represents a TCP server
type TcpServer struct {
	ctx *svc.ServerCtx
}

// NewTcpServer creates a new TCP server instance
func NewTcpServer(ctx *svc.ServerCtx) *TcpServer {
	return &TcpServer{ctx: ctx}
}

// Listen starts the TCP server and listens for incoming connections
func (s *TcpServer) Listen() error {
	ln, err := net.Listen("tcp", s.ctx.Config.ListenOn)
	if err != nil {
		log.Error().Err(err).Msg("Error starting server")
		return err
	}
	defer ln.Close()

	log.Info().Msgf("Listening on %s", s.ctx.Config.ListenOn)
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Error accepting connection")
			continue
		}
		go s.handleConn(conn)
	}
}

// handleConn manages an individual connection from a client
func (s *TcpServer) handleConn(conn net.Conn) {
	log.Info().Msg("Connection established with client")

	reader := bufio.NewReader(conn)
	for {
		buf := make([]byte, 2048)
		n, err := reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				log.Info().Msg("Connection closed by client")
			} else {
				log.Error().Err(err).Msg("Error reading from client")
			}
			break
		}
		s.processMessage(buf[:n], conn)
	}
}

// processMessage processes incoming messages from the client
func (s *TcpServer) processMessage(messageStr []byte, conn net.Conn) {
	log.Debug().Msgf("Processing message")
	m, err := message.Unmarshal(messageStr)
	if err != nil {
		// 从外部接收到的数据，转发到内部
		messageId := utils.GenerateID()
		log.Debug().Msgf("New message ID: %s", messageId)
		s.ctx.Messages[messageId] = conn
		s.handleData(message.Message{
			Type: message.MessageTypeData,
			Data: messageStr,
			Id:   messageId,
		})
		return
	}

	switch m.Type {
	case message.MessageTypeConnect:
		s.HandleConnect(*m, conn)
	case message.MessageTypeData:
		messageIds := make([]string, 0)
		for id := range s.ctx.Messages {
			messageIds = append(messageIds, id)
		}
		log.Debug().Msgf("Messages %v", messageIds)
		s.ctx.Messages[m.Id].Write(m.Data)
		s.ctx.Messages[m.Id].Close()
		delete(s.ctx.Messages, m.Id)
	case message.MessageTypeDisconnect:
		log.Info().Msg("Client disconnected")
	case message.MessageTypeHeartbeat:
		s.sendHeartbeatResponse(conn)
	default:
		log.Warn().Msg("Unknown message type")
	}
}

// handleConnect manages the connection request from the client
func (s *TcpServer) HandleConnect(m message.Message, conn net.Conn) {
	clientConfig, err := config.ParseFromEncoded(string(m.Data))
	if err != nil {
		log.Error().Err(err).Msg("Error parsing client config")
		s.sendDisconnectResponse(conn, "Invalid config")
		conn.Close()
		return
	}

	tunnel, err := s.ctx.TunnelModel.GetTunnelByID(clientConfig.TunnelID)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving tunnel")
		s.sendDisconnectResponse(conn, "Tunnel not found")
		conn.Close()
		return
	}

	tunnel.Status = "online"
	s.ctx.TunnelModel.Update(tunnel)
	s.ctx.Tunnels[tunnel.ID] = conn

	response := message.Message{
		Type: message.MessageTypeConnect,
		Data: []byte(fmt.Sprintf("Connected to tunnel %s", clientConfig.TunnelID)),
	}
	if err := s.sendResponse(conn, response); err != nil {
		log.Error().Err(err).Msg("Failed to send response")
	}
	log.Info().Msgf("Client connected: %s", clientConfig.TunnelID)
}

// handleData processes the data from the client
func (s *TcpServer) handleData(m message.Message) {
	tunnelID := ""
	// todo: 消息规则知道转发到哪个隧道
	// 通过隧道ID获取隧道连接
	if utils.HttpPattern.Match(m.Data) {
		// todo: 处理消息 把消息host 转换为规则host
		host := utils.GetHostFromHttpMessage(m.Data)
		for _, route := range s.ctx.Routes {
			if route.Hostname == host {
				m.Protocol = route.Protocol
				tunnelID = route.TunnelID
				m.Target = route.Target
				break
			}
		}
	}

	if utils.SshPattern.Match(m.Data) {
		m.Protocol = model.TypeSsh
	}

	// 通过隧道ID获取隧道连接
	if _, ok := s.ctx.Tunnels[tunnelID]; !ok {
		log.Error().Msg("tunnel not found")
		return
	}

	conn := s.ctx.Tunnels[tunnelID]
	mData, err := m.Marshal()
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling message")
		return
	}
	_, err = conn.Write(mData)
	if err != nil {
		log.Error().Err(err).Msg("Error sending message")
		return
	}
	log.Debug().Msg("Data sent to tunnel")
}

// sendHeartbeatResponse sends a heartbeat response to the client
func (s *TcpServer) sendHeartbeatResponse(conn net.Conn) {
	response := message.Message{
		Type: message.MessageTypeHeartbeat,
		Data: []byte("pong"),
	}
	if err := s.sendResponse(conn, response); err != nil {
		log.Error().Err(err).Msg("Failed to send heartbeat response")
	}
}

// sendResponse sends a message back to the client
func (s *TcpServer) sendResponse(conn net.Conn, response message.Message) error {
	data, err := response.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	if _, err = conn.Write(data); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}
	return nil
}

// sendDisconnectResponse sends a disconnect message to the client
func (s *TcpServer) sendDisconnectResponse(conn net.Conn, reason string) {
	response := message.Message{
		Type: message.MessageTypeDisconnect,
		Data: []byte(reason),
	}
	s.sendResponse(conn, response)
}
