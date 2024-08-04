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
	DefaultListenOn = ":5429"
	DefaultApi      = ":8000"
)

var (
	clients = make(map[string]net.Conn)
	channel = make(chan string)
)

func main() {
	listenOn := os.Getenv("ListenOn")
	api := os.Getenv("Api")

	if listenOn == "" {
		listenOn = DefaultListenOn
	}
	if api == "" {
		api = DefaultApi
	}
	db, err := gorm.Open(sqlite.Open("tunnel.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&model.Tunnel{})
	lister, err := net.Listen("tcp", listenOn)
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
