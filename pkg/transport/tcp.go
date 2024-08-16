package transport

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/errors"
	"github.com/obud-dev/tunnel/pkg/message"
	"github.com/obud-dev/tunnel/pkg/model"
	"github.com/obud-dev/tunnel/pkg/svc"
	"github.com/obud-dev/tunnel/pkg/utils"
)

// TCP Client Constants
const (
	heartbeatInterval = 60 * time.Second // 心跳包发送间隔
	reconnectAttempts = 3                // 最大重连尝试次数
	reconnectInterval = 10 * time.Second // 每次重连间隔
	heartTimeout      = 5 * time.Second  // 心跳包超时时间
)

// TcpClient is a client for the TCP protocol
type TcpClient struct {
	conn    net.Conn
	conf    *config.ClientConfig
	channel chan []byte
}

// NewTcpClient creates a new TCP client
func NewTcpClient(token string) (*TcpClient, error) {
	conf, err := config.ParseFromEncoded(token)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &TcpClient{conf: conf, channel: make(chan []byte)}, nil
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
	m := message.Message{
		Id:   utils.GenerateID(),
		Type: message.MessageTypeConnect,
		Data: []byte(c.conf.TunnelID),
	}
	req, err := m.Marshal()
	if err != nil {
		return fmt.Errorf("failed to marshal connect message: %w", err)
	}
	c.conn.Write(req)

	go c.readLoop()
	go c.Heartbeat()
	go c.sendToServer()
	select {}
}

// readLoop handles incoming messages from the server
func (c *TcpClient) readLoop() {
	defer c.conn.Close()

	reader := bufio.NewReader(c.conn)
	for {
		data, err := reader.ReadBytes('}')
		if err != nil {
			if err == io.EOF {
				log.Info().Msg("Connection closed by server")
			} else {
				log.Error().Err(err).Msg("Error reading from server")
			}
			break
		}

		m, err := message.Unmarshal(data)
		if err != nil {
			log.Error().Err(err).Msgf("Failed to unmarshal message %s", data)
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

// SendMessage encrypt data and ready to send
func (c *TcpClient) SendMessage(m message.Message) error {
	data, err := m.Encrypt(c.conf.Token)
	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message")
		return err
	}
	c.channel <- data
	return nil
}

// sendToServer send data to server
func (c *TcpClient) sendToServer() {
	defer close(c.channel)

	for message := range c.channel {
		_, err := c.conn.Write(message)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message to server")
			return
		}
		log.Debug().Msg("Data sent to server")
	}
}

// RecieveData processes data received from the server
func (c *TcpClient) RecieveData(m message.Message) {
	log.Debug().Msg("Received data from server")

	// Decrypt data
	data, err := m.Decrypt(c.conf.Token)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decrypt message")
		return
	}
	m.Data = data

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
	timer := time.NewTimer(heartTimeout)
	defer timer.Stop()

	for {
		time.Sleep(heartbeatInterval)
		done := make(chan error, 1)
		go func() {
			m := message.Message{
				Id:   utils.GenerateID(),
				Type: message.MessageTypeHeartbeat,
				Data: []byte("ping"),
			}
			err := c.SendMessage(m)
			done <- err
		}()
		if !timer.Stop() {
			// 如果计时器已经触发，清空通道
			<-timer.C
		}
		timer.Reset(heartTimeout)
		select {
		case err := <-done:
			if err != nil {
				c.ReconnectToServer()
				return
			}
		case <-timer.C:
			fmt.Print("Send heartbeat timeout")
			c.ReconnectToServer()
			return
		}
	}
}

// ReconnectToServer attempts to reconnect to the TCP server
func (c *TcpClient) ReconnectToServer() error {
	if c.conn != nil {
		c.conn.Close()
	}
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
	headBuf := make([]byte, 1)
	for {
		_, err := reader.Read(headBuf)
		if err != nil {
			break
		}
		if strings.Contains(string(headBuf), "{") {
			data, err := reader.ReadBytes('}')
			if err != nil {
				if err == io.EOF {
					log.Info().Msg("Connection closed by client")
				} else {
					log.Error().Err(err).Msg("Error reading from client")
				}
				break
			}
			data = append(headBuf, data...)
			s.processMessage(data, conn)
			continue
		}
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
		data := append(headBuf, buf[:n]...)
		s.processMessage(data, conn)
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
		tunnel := &svc.ActiveTunnel{
			Conn:    conn,
			Channel: make(chan []byte),
		}
		// map 不是并发安全的，阻止其它线程操作map，操作即时完成，不影响性能（应该）
		s.ctx.Mutex.Lock()
		s.ctx.Messages[messageId] = tunnel
		s.ctx.Mutex.Unlock()
		m := message.Message{
			Type: message.MessageTypeData,
			Data: messageStr,
			Id:   messageId,
		}
		s.handleData(m)
		go s.sendToVisitor(messageId)
		return
	}

	switch m.Type {
	case message.MessageTypeConnect:
		s.HandleConnect(*m, conn)
	case message.MessageTypeData:
		go s.handleClientData(*m)
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
	TunnelID := string(m.Data)

	tunnel, err := s.ctx.TunnelModel.GetTunnelByID(TunnelID)
	if err != nil {
		log.Error().Err(err).Msg("Error retrieving tunnel")
		s.sendDisconnectResponse(conn, "Tunnel not found")
		conn.Close()
		return
	}

	tunnel.Status = "online"
	s.ctx.TunnelModel.Update(tunnel)
	s.ctx.Tunnels[tunnel.ID] = &svc.ActiveTunnel{
		Conn:    conn,
		Token:   tunnel.Token,
		Channel: make(chan []byte),
	}
	go s.sendToClient(TunnelID)
	response := message.Message{
		Id:   m.Id,
		Type: message.MessageTypeConnect,
		Data: []byte(fmt.Sprintf("Connected to tunnel %s", TunnelID)),
	}
	if err := s.sendResponse(conn, response); err != nil {
		log.Error().Err(err).Msg("Failed to send response")
	}
	log.Info().Msgf("Client connected: %s", TunnelID)
}

// handleData processes the data from the visitor
func (s *TcpServer) handleData(m message.Message) {
	tunnelID := ""
	// 消息规则知道转发到哪个隧道
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
		s.ctx.Messages[m.Id].Channel <- []byte(errors.Http500)
		return
	}

	tunnel := s.ctx.Tunnels[tunnelID]
	s.ctx.Messages[m.Id].Token = tunnel.Token
	mData, err := m.Encrypt(tunnel.Token)
	if err != nil {
		log.Error().Err(err).Msg("Error marshalling message")
		s.ctx.Messages[m.Id].Channel <- []byte(errors.Http500)
		return
	}
	tunnel.Channel <- mData
}

// handleClientData processes the data from the client
func (s *TcpServer) handleClientData(m message.Message) {
	if _, ok := s.ctx.Messages[m.Id]; !ok {
		log.Error().Msg("message conn not found")
		return
	}
	tunnel := s.ctx.Messages[m.Id]
	data, err := m.Decrypt(tunnel.Token)
	if err != nil {
		fmt.Println("Error to decrypt data:", err)
		return
	}
	tunnel.Channel <- data
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

// sendToVisitor sends client response to visitor with deal all Producter-Goroutine
func (s *TcpServer) sendToVisitor(mid string) {
	m := s.ctx.Messages[mid]
	defer close(m.Channel)
	defer delete(s.ctx.Messages, mid)

	for message := range m.Channel {
		_, err := m.Conn.Write(message)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message")
			return
		}
		log.Debug().Msg("Data sent to visitor")
	}
	m.Conn.Close()
}

// sendToClient sends to client with deal all Producter-Goroutine
func (s *TcpServer) sendToClient(tid string) {
	tunnel := s.ctx.Tunnels[tid]
	defer close(tunnel.Channel)

	for message := range tunnel.Channel {
		_, err := tunnel.Conn.Write(message)
		if err != nil {
			log.Error().Err(err).Msg("Error sending message")
			return
		}
		log.Debug().Msg("Data sent to tunnel")
	}
}
