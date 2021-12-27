package xconsul

import (
	"github.com/HaleyLeoZhang/go-component/helper"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

func NewClient(cfg *Config) (output *Client, err error) {
	// 创建连接consul服务配置
	config := consulapi.DefaultConfig()
	config.Address = cfg.Addr
	client, err := consulapi.NewClient(config)
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	output = &Client{
		clt:      client,
		register: NewRegister(cfg.ServiceName, helper.GetLocalIpV4(), cfg.HttpPort, cfg.HealthCheckRouter),
	}
	return
}
