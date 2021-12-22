package xconsul

import (
	"github.com/google/uuid"
	consulapi "github.com/hashicorp/consul/api"
)

type Client struct {
	clt      *consulapi.Client
	register *Register
}

type Config struct {
	Addr string `yaml:"addr"`
}

type Register struct {
	// 服务注册信息
	UniqueName string // 一般是 uuid
	Name       string
	// 当前机器的 网络信息
	Ip     string // 本机IP
	Port   int
	Router string // 健康检查的路由如 /health
	// 健康检查
	HealthTimeout                  string // 如 3s
	HealthInterval                 string // 如 3s
	DeregisterCriticalServiceAfter string // 故障检查失败指定秒数后 consul自动将注册服务删除 如 30s
}

func NewRegister(serviceName string, ip string, port int, router string) (reg *Register) {
	reg = &Register{
		UniqueName: uuid.New().String(),
		Name:       serviceName,
		Ip:         ip,
		Port:       port,
		Router:     router,
		// 写死
		HealthTimeout:                  "3s",
		HealthInterval:                 "3s",
		DeregisterCriticalServiceAfter: "30s",
	}
	return
}
