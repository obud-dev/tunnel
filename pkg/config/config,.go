package config

type ClientConfig struct {
	TunnelID string `json:"tunnel_id"` // 隧道 ID
	Token    string `json:"token"`     // 隧道 令牌
	Server   string `json:"server"`    // 服务器地址
}

type ServerConfig struct {
	Addr string `json:"addr"` // 服务器地址
}
