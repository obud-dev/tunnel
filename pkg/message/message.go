package message

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"

	"github.com/obud-dev/tunnel/pkg/model"
)

type MessageType int

const (
	MessageTypeData MessageType = iota
	MessageTypeConnect
	MessageTypeDisconnect
	MessageTypeHeartbeat
)

type Message struct {
	Type     MessageType    `json:"type"`
	Data     []byte         `json:"data"`
	Id       string         `json:"id"`
	Protocol model.Protocol `json:"protocol"`
	Target   string         `json:"target"`
}

func (m *Message) Marshal() ([]byte, error) {
	return json.Marshal(m)
}

func Unmarshal(data []byte) (*Message, error) {
	m := &Message{}
	err := json.Unmarshal(data, m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// 使用AES-GCM加解密
func (m *Message) Encrypt(key string) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// 使用AES-GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// 生成随机的nonce
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	// 加密并附加认证标签
	ciphertext := aesGCM.Seal(nonce, nonce, m.Data, nil)

	m.Data = ciphertext
	return m.Marshal()
}

func (m *Message) Decrypt(key string) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// 使用AES-GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	// 提取nonce
	nonceSize := aesGCM.NonceSize()
	nonce, ciphertext := m.Data[:nonceSize], m.Data[nonceSize:]
	// 解密
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("token authentication failed or message has falsified")
	}

	return plaintext, nil
}
