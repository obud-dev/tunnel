package token

import (
	"encoding/base64"
	"encoding/json"

	"github.com/google/uuid"
)

type Token struct {
	UID      string `json:"uid"`
	TunnelID string `json:"tunnel_id"`
	Token    string `json:"token"`
}

// 生成 Token
func GenerateToken(uid string) (string, error) {
	t := Token{UID: uid, TunnelID: uuid.New().String()}
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
