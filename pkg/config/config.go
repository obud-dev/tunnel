package config

import (
	"encoding/base64"
	"encoding/json"
)

type ClientConfig struct {
	TunnelID string `json:"tunnel_id"` // 隧道 ID
	Token    string `json:"token"`     // 隧道 令牌,密钥
	Server   string `json:"server"`    // 服务器地址
}

type ServerConfig struct {
	Host     string `json:"host"`      // 服务器地址
	ListenOn string `json:"listen_on"` // 监听地址 (默认 :5429)
	Api      string `json:"api"`       // API地址 (默认 :8000)
	Domain   string `json:"domain"`    // 域名 生成client token时使用
	User     string `json:"user"`      // api 用户名
	Password string `json:"password"`  // api 密码
}

func ParseFromEncoded(encoded string) (*ClientConfig, error) {
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}
	config := &ClientConfig{}
	err = json.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func (c *ClientConfig) Encode() (string, error) {
	jsonData, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(jsonData), nil
}
