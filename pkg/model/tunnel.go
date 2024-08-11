package model

import "gorm.io/gorm"

type Tunnel struct {
	ID     string `json:"id" gorm:"primaryKey"`
	Name   string `json:"name" gorm:"unique"`
	Status string `json:"status" gorm:"default:'offline'"`
	Uptime int64  `json:"uptime"`
	Token  string `json:"token" gorm:"not null"` // 内网进程用来连接公网服务tunnel的token
}

func (t *Tunnel) TableName() string {
	return "tunnels"
}

type defaultTunnelModel struct {
	db *gorm.DB
}

type TunnelModel interface {
	GetTunnels() ([]Tunnel, error)
	GetTunnelByID(id string) (*Tunnel, error)
	Insert(tunnel *Tunnel) error
	Update(tunnel *Tunnel) error
	Delete(tunnel *Tunnel) error
}

func NewTunnelModel(db *gorm.DB) *defaultTunnelModel {
	return &defaultTunnelModel{db: db}
}

func (m *defaultTunnelModel) GetTunnels() ([]Tunnel, error) {
	var tunnels []Tunnel
	err := m.db.Find(&tunnels).Error
	if err != nil {
		return nil, err
	}
	return tunnels, nil
}

func (m *defaultTunnelModel) GetTunnelByID(id string) (*Tunnel, error) {
	var tunnel Tunnel
	err := m.db.First(&tunnel, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &tunnel, nil
}

func (m *defaultTunnelModel) Insert(tunnel *Tunnel) error {
	return m.db.Create(tunnel).Error
}

func (m *defaultTunnelModel) Update(tunnel *Tunnel) error {
	return m.db.Save(tunnel).Error
}

func (m *defaultTunnelModel) Delete(tunnel *Tunnel) error {
	return m.db.Delete(tunnel).Error
}
