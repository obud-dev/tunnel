package transport

// 使用tcp协议进行通信

import (
	"net"

	"github.com/obud-dev/tunnel/pkg/config"
)

// TCPClient is a client for the TCP protocol
type TCPClient struct {
	conn   net.Conn
	config *config.ClientConfig
}

// NewTCPClient creates a new TCP client

type TCPServer struct {
	listener net.Listener
	config   *config.ServerConfig
}

// NewTCPServer creates a new TCP server
func NewTCPServer(config *config.ServerConfig) (*TCPServer, error) {
	listener, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, err
	}
	return &TCPServer{listener: listener, config: config}, nil
}

// Close closes the server
func (s *TCPServer) Close() error {
	return s.listener.Close()
}

// Accept waits for a connection
func (s *TCPServer) Accept() (net.Conn, error) {
	return s.listener.Accept()
}

// NewTCPClient creates a new TCP client
func NewTCPClient(config *config.ClientConfig) (*TCPClient, error) {
	conn, err := net.Dial("tcp", config.Server)
	if err != nil {
		return nil, err
	}
	return &TCPClient{conn: conn, config: config}, nil
}

// Close closes the client
func (c *TCPClient) Close() error {
	return c.conn.Close()
}

// Read reads data from the connection
func (c *TCPClient) Read(b []byte) (int, error) {
	return c.conn.Read(b)
}

// Write writes data to the connection
func (c *TCPClient) Write(b []byte) (int, error) {
	return c.conn.Write(b)
}
