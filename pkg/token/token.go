package token

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

type Token struct {
	TunnelID string `json:"tunnel_id"` // 隧道 ID
	Token    string `json:"token"`     // 隧道 令牌
	Server   string `json:"server"`    // 服务器地址
}

// 生成 Token
func GenerateToken(uid, server string) (string, error) {
	t := Token{TunnelID: uuid.New().String(), Server: server}
	jsonData, err := json.Marshal(t)
	if err != nil {
		return "", err
	}
	encodedToken := base64.StdEncoding.EncodeToString(jsonData)
	return encodedToken, nil
}

// 解析 Token
func ParseToken(encodedToken string) (*Token, error) {
	jsonData, err := base64.StdEncoding.DecodeString(encodedToken)
	if err != nil {
		return nil, err
	}
	var t Token
	err = json.Unmarshal(jsonData, &t)
	if err != nil {
		return nil, err
	}
	return &t, nil
}
