package config

import (
	"encoding/base64"
	"encoding/json"
)

type ClientConfig struct {
	TunnelID string `json:"tunnel_id"` // 隧道 ID
	Token    string `json:"token"`     // 隧道 令牌
	Server   string `json:"server"`    // 服务器地址
}

type ServerConfig struct {
	TunnelAddr string `json:"tunnel_addr"` // 隧道地址
	ServerAddr string `json:"server_addr"` // 服务器地址
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
