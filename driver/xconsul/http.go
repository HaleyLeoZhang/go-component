package xconsul

// http 注册模式

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
)

// 注册服务 --- 服务启动时

func (c *Client) HttpRegister() (err error) {
	// 创建注册到consul的服务到
	registration := new(consulapi.AgentServiceRegistration)
	registration.ID = c.register.UniqueName
	registration.Name = c.register.Name // 根据这个名称来找这个服务
	registration.Port = c.register.Port
	//registration.Tags = []string{"shitingbao_test_service"} //这个就是一个标签，可以根据这个来找这个服务，相当于V1.1这种
	registration.Address = c.register.Ip

	// 增加consul健康检查回调函数
	check := new(consulapi.AgentServiceCheck)
	check.HTTP = fmt.Sprintf("http://%s:%d%v", registration.Address, registration.Port, c.register.Router)
	check.Timeout = c.register.HealthTimeout
	check.Interval = c.register.HealthInterval
	check.DeregisterCriticalServiceAfter = c.register.DeregisterCriticalServiceAfter
	registration.Check = check

	// 注册服务到consul
	err = c.clt.Agent().ServiceRegister(registration)
	if err != nil {
		return
	}
	println(fmt.Sprintf("Registered to consul %v:%v  %v", registration.Address, registration.Port, c.register.UniqueName))
	return
}

// 注销服务 --- 监听服务关闭 -- 无论 HTTP/RPC

func (c *Client) Deregister() (err error) {
	err = c.clt.Agent().ServiceDeregister(c.register.UniqueName)
	if err != nil {
		return
	}
	println(fmt.Sprintf("Deregister to consul %v", c.register.UniqueName))
	return
}
