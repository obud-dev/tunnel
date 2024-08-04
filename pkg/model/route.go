package model

import "gorm.io/gorm"

// 在服务端创建，用于将公网请求转发到内网服务。并内网使用转发

type Protocol string

const (
	TypeHttp Protocol = "http"
	TypeTcp  Protocol = "tcp"
	TypeUdp  Protocol = "udp"
	TypeSsh  Protocol = "ssh"
	TypeRdp  Protocol = "rdp"
)

type Route struct {
	ID       string   `gorm:"primaryKey"`
	TunnelID string   `gorm:"not null"` // 路由所属的隧道
	Hostname string   `gorm:"not null"` // 域名
	Prefix   string   // 路由前缀
	Target   string   `gorm:"not null"` // 内网目标服务地址
	Protocol Protocol `gorm:"not null"` // 协议
}

func (r *Route) TableName() string {
	return "routes"
}

type defaultRouteModel struct {
	db *gorm.DB
}

type RouteModel interface {
	GetRoutes() ([]Route, error)
	GetRouteByID(id string) (*Route, error)
	GetRoutesByTunnelID(tunnelID string) ([]Route, error)
	Insert(route *Route) error
	Update(route *Route) error
	Delete(route *Route) error
}

func NewRouteModel(db *gorm.DB) *defaultRouteModel {
	return &defaultRouteModel{db: db}
}

func (m *defaultRouteModel) GetRoutes() ([]Route, error) {
	var routes []Route
	err := m.db.Find(&routes).Error
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (m *defaultRouteModel) GetRouteByID(id string) (*Route, error) {
	var route Route
	err := m.db.First(&route, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &route, nil
}

func (m *defaultRouteModel) GetRoutesByTunnelID(tunnelID string) ([]Route, error) {
	var routes []Route
	err := m.db.Where("tunnel_id = ?", tunnelID).Find(&routes).Error
	if err != nil {
		return nil, err
	}
	return routes, nil
}

func (m *defaultRouteModel) Insert(route *Route) error {
	return m.db.Create(route).Error
}

func (m *defaultRouteModel) Update(route *Route) error {
	return m.db.Save(route).Error
}

func (m *defaultRouteModel) Delete(route *Route) error {
	return m.db.Delete(route).Error
}
