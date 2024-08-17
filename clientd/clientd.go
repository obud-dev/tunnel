package main

import (
	"flag"

	"github.com/obud-dev/tunnel/pkg/transport"
	"github.com/obud-dev/tunnel/pkg/utils"
	"github.com/rs/zerolog/log"
)

func main() {
	token := flag.String("token", "", "token to connect to the server")
	flag.Parse()

	if *token == "" {
		log.Error().Msg("token is required")
		return
	}

	// 打印内存使用情况
	go utils.PrintMemoryUsage()

	// 连接到公网服务器
	client, err := transport.NewTcpClient(*token)
	if err != nil {
		log.Error().Err(err).Msg("failed to create client")
		return
	}
	client.Connect()
}
