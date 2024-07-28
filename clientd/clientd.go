package main

import (
	"fmt"
	"log"
	"net"
)

const (
	remoteAddr = "0.0.0.0:5429"
	tunnelId   = "8b2526ad-aef1-47da-85c5-ef66c4666642"
	localAddr  = "0.0.0.0:8080"
)

func main() {

	// 连接到公网服务器
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer conn.Close()

	// todo: 这里可以实现身份验证、注册等功能
	json := []byte(`{"tunnel_id":"` + tunnelId + `"}`)
	conn.Write(json)

	// 监听来自公网服务器的流量
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println("Error reading from server:", err)
			return
		}

		fmt.Println("read:", string(buffer[:n]))
		// 将接收到的数据转发到内网服务
		// handleRequest(buffer[:n], conn)
	}

}
func handleRequest(data []byte, conn net.Conn) {
	// 连接内网服务
	serviceConn, err := net.Dial("tcp", localAddr)
	if err != nil {
		log.Println("Failed to connect to service:", err)
		return
	}
	defer serviceConn.Close()

	// 转发数据到内网服务
	_, err = serviceConn.Write(data)
	if err != nil {
		log.Println("Failed to write to service:", err)
	}

	// 读取内网服务的响应
	responseBuffer := make([]byte, 1024)
	n, err := serviceConn.Read(responseBuffer)
	if err != nil {
		log.Println("Error reading from service:", err)
		return
	}

	// 将内网服务的响应转发到公网服务器
	_, err = conn.Write(responseBuffer[:n])
	if err != nil {
		log.Println("Failed to write to server:", err)
	}
}
