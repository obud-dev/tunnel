package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/obud-dev/tunnel/pkg/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const (
	// DefaultBindAddr is the default address to bind to
	DefaultBindAddr = ":5429"
	// DefaultServerAddr is the default address of the server
	DefaultServerAddr = ":8000"
)

var (
	clients = make(map[string]net.Conn)
	channel = make(chan string)
)

func main() {
	bind_addr := os.Getenv("BIND_ADDR")
	server_addr := os.Getenv("SERVER_ADDR")

	if bind_addr == "" {
		bind_addr = DefaultBindAddr
	}
	if server_addr == "" {
		server_addr = DefaultServerAddr
	}
	db, err := gorm.Open(sqlite.Open("tunnel.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Tunnel{})
	lister, err := net.Listen("tcp", bind_addr)
	if err != nil {
		panic(err)
	}
	for {
		conn, err := lister.Accept()
		if err != nil {
			panic(err)
		}
		go func() {
			fmt.Println("new connection")
			// defer conn.Close()
			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				fmt.Println("error reading:", err)
				return
			}
			fmt.Println("read:", string(buf[:n]))
			// 收到 {"tunnel_id":"8b2526ad-aef1-47da-85c5-ef66c4666642"}
			// 从中提取tunnel_id
			jsonData := make(map[string]string)
			err = json.Unmarshal(buf[:n], &jsonData)
			if err != nil {
				fmt.Println("error unmarshal:", err)
			}
			tunnelId := jsonData["tunnel_id"]
			fmt.Println("tunnel_id:", tunnelId)
			// 从数据库中查找tunnel_id对应的记录
			// var tunnel model.Tunnel
			// db.First(&tunnel, "id = ?", tunnelId)
			// if tunnel.ID == "" {
			// 	fmt.Println("tunnel not found")
			// 	return
			// }
			// 连接到公网服务器
			clients[tunnelId] = conn
			conn.Write([]byte("connected"))
			fmt.Println("connected")
		}()
	}

}
