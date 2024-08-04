package utils

import (
	"bytes"
	"errors"
	"fmt"
	"io"
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

func ResponseToBytes(resp *http.Response) ([]byte, error) {
	// 创建一个 bytes.Buffer 用于拼接响应内容
	var buf bytes.Buffer

	// 写入状态行
	buf.WriteString(fmt.Sprintf("HTTP/%d.%d %s\r\n", resp.ProtoMajor, resp.ProtoMinor, resp.Status))

	// 写入头部
	for key, values := range resp.Header {
		for _, value := range values {
			buf.WriteString(fmt.Sprintf("%s: %s\r\n", key, value))
		}
	}
	buf.WriteString("\r\n") // 头部与主体之间的空行

	// 读取主体内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// 写入主体内容
	buf.Write(body)

	// 关闭响应主体
	// 注意：一旦读取了主体，Resp.Body 将被关闭。
	// 因此，如果需要使用响应体，请考虑先将其保存。
	// resp.Body.Close()  // 这里不关闭，因为我们还需要它的内容

	return buf.Bytes(), nil
}
