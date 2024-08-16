package main

import (
	"github.com/obud-dev/tunnel/pkg/transport"
	"github.com/obud-dev/tunnel/pkg/utils"
	"github.com/rs/zerolog/log"
)

const (
	token = "eyJ0dW5uZWxfaWQiOiJjY2Y3MjU4Zi0wZTQxLTRlODAtYTRlYS0xOGVkODE5NWI5OGUiLCJ0b2tlbiI6IjEyMzRhYmNkNTY3OGVmOTAxMjM0YWJjZCIsInNlcnZlciI6IjAuMC4wLjA6NTQyOSJ9"
)

func main() {

	go utils.PrintMemoryUsage()

	// 连接到公网服务器
	client, err := transport.NewTcpClient(token)
	if err != nil {
		log.Error().Err(err).Msg("failed to create client")
		return
	}
	client.Connect()
}
