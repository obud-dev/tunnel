package model

type Tunnel struct {
	ID     string `gorm:"primaryKey"`
	Name   string `gorm:"unique"`
	Status string `gorm:"default:'offline'"`
	Uptime int64
	Token  string `gorm:"not null"` // 内网进程用来连接公网服务tunnel的token
}

func (t *Tunnel) TableName() string {
	return "tunnels"
}
