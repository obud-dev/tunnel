package main

import (
	"flag"
	"net/http"
	"sync"
	"time"

	"github.com/obud-dev/tunnel/pkg/utils"
	"github.com/rs/zerolog/log"
)

var client = &http.Client{
	Timeout: time.Second * 10, // 设置超时时间
}

// makeRequest 发起 HTTP GET 请求并打印请求时间和响应状态
func makeRequest(url string, t, n int) {
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		log.Error().Err(err).Msgf("Failed to make request to %s", url)
		return
	}
	defer resp.Body.Close()
	log.Info().Msgf("threads: %d, request: %d, status: %s, duration: %s", t, n, resp.Status, duration)
}

func main() {
	// 使用 flag 包定义命令行参数
	url := flag.String("url", "", "URL to test (required)")
	t := flag.Int("t", 1, "Number of concurrent threads")
	n := flag.Int("n", 1, "Number of requests per thread")

	// 解析命令行参数
	flag.Parse()

	// 验证 URL 是否被提供
	if *url == "" {
		flag.Usage()
		return
	}

	utils.InitLogger()
	go utils.PrintMemoryUsage()

	// 使用 WaitGroup 等待所有请求完成
	var wg sync.WaitGroup

	// 启动指定数量的线程
	for i := 0; i < *t; i++ {
		wg.Add(1)
		go func(t int) {
			defer wg.Done()
			// 每个线程发起指定数量的请求
			for j := 0; j < *n; j++ {
				makeRequest(*url, t, j)
			}
		}(i)
	}

	// 等待所有线程完成
	wg.Wait()
}
