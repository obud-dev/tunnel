package utils

import (
	"errors"
	"fmt"
	"net"
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
