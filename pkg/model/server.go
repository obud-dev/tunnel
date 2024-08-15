package model

import "gorm.io/gorm"

type Server struct {
	Host      string `json:"host" gorm:"primaryKey"` // 主机 IP地址
	ListenOn  string `json:"listen_on"`              // 监听端口(默认 :5429)
	Api       string `json:"api"`                    // WEB_UI/API端口(默认 :8000)
	Domain    string `json:"domain"`                 // 域名 (可选) 如果没有设置则为 Host:ListenOn
	ApiDomain string `json:"api_domain"`             // WEB_UI/API域名 (可选) 如果没有设置则为 Host:Api
	Version   string `json:"version"`                // 版本
}

func (s *Server) TableName() string {
	return "servers"
}

type defaultServerModel struct {
	db *gorm.DB
}

type ServerModel interface {
	GetServer() (*Server, error)
	Update(server *Server) error
}

func NewServerModel(db *gorm.DB) *defaultServerModel {
	return &defaultServerModel{db: db}
}

func (m *defaultServerModel) GetServer() (*Server, error) {
	var server Server
	err := m.db.First(&server).Error
	if err != nil {
		return nil, err
	}
	return &server, nil
}

func (m *defaultServerModel) Update(server *Server) error {
	return m.db.Save(server).Error
}
