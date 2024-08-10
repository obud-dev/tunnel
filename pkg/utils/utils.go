package utils

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/google/uuid"
)

func GetAvailablePort(min int) (int, error) {
	// 获取一个可用的端口
	for port := min; port < 65535; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
		if err == nil {
			ln.Close()
			return port, nil
		}
	}
	return 0, errors.New("no available port")
}

func GenerateID() string {
	// 生成一个随机ID
	return uuid.New().String()
}

// GetHostFromHttpMessage extracts the host from the HTTP message
func GetHostFromHttpMessage(m []byte) string {
	reader := bufio.NewReader(bytes.NewReader(m))
	req, err := http.ReadRequest(reader)
	if err != nil {
		return ""
	}
	return req.Host
}
