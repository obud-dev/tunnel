package main

import (
	"flag"
	"net/http"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// makeRequest 发起 HTTP GET 请求并打印请求时间和响应状态
func makeRequest(url string, t, n int) {
	start := time.Now()
	resp, err := http.Get(url)
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

	// 打印内存使用情况
	go func() {
		for {
			printMemoryUsage()
			time.Sleep(5 * time.Second)
		}
	}()

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

func printMemoryUsage() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.Info().Msgf("Alloc = %v MiB TotalAlloc = %v MiB Sys = %v MiB NumGC = %v", bToMb(m.Alloc), bToMb(m.TotalAlloc), bToMb(m.Sys), m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
