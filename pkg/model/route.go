package model

// 在服务端创建，用于将公网请求转发到内网服务。并内网使用转发
type Route struct {
	ID       string `gorm:"primaryKey"`
	TunnelID string `gorm:"not null"` // 路由所属的隧道
	Hostname string `gorm:"not null"` // 域名
	Prefix   string // 路由前缀
	Target   string `gorm:"not null"` // 内网目标服务地址
	Type     string `gorm:"not null"` // 转发类型 http tcp udp ssh rdp
}

func (r *Route) TableName() string {
	return "routes"
}
