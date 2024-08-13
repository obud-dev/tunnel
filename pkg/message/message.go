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

// 使用AES-CTR加解密
func (m *Message) Encrypt(key string) ([]byte, error) {
	if len(key) != 16 && len(key) != 24 && len(key) != 32 {
		return nil, fmt.Errorf("invalid key size: must be 16, 24, or 32 bytes")
	}
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}
	// 创建一个随机的IV
	ciphertext := make([]byte, aes.BlockSize+len(m.Data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	// 加密
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], m.Data)
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

	iv := m.Data[:aes.BlockSize]
	ciphertext := m.Data[aes.BlockSize:]

	// 解密
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}
