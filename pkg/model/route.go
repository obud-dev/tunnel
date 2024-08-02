package model

// 在服务端创建，用于将公网请求转发到内网服务。并内网使用转发

type RouteType string

const (
	RouteTypeHttp RouteType = "http"
	RouteTypeTcp  RouteType = "tcp"
	RouteTypeUdp  RouteType = "udp"
	RouteTypeSsh  RouteType = "ssh"
	RouteTypeRdp  RouteType = "rdp"
)

type Route struct {
	ID       string    `gorm:"primaryKey"`
	TunnelID string    `gorm:"not null"` // 路由所属的隧道
	Hostname string    `gorm:"not null"` // 域名
	Prefix   string    // 路由前缀
	Target   string    `gorm:"not null"` // 内网目标服务地址
	Type     RouteType `gorm:"not null"` // 转发类型
}

func (r *Route) TableName() string {
	return "routes"
}
