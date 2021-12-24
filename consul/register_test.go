package consul

import (
	"testing"

	"github.com/hashicorp/consul/api"
)

// go注册服务到consul
func TestRegistry(t *testing.T) {
	// 1. 配置连接
	cfg := api.DefaultConfig()
	cfg.Address = "192.168.5.130:8500"

	// 2. 创建连接
	client, err := api.NewClient(cfg)
	if err != nil {
		panic("create client failed:" + err.Error())
	}

	// 3. 创建健康检查
	check := &api.AgentServiceCheck{
		HTTP:                           "http://10.225.20.217:8080/health",
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "10s",
	}

	// 4. 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Name = "service-name"
	registration.ID = "service-name"
	registration.Port = 8080
	registration.Tags = []string{"service", "mooc"}
	registration.Address = "10.225.20.217"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil {
		panic("naming service failed:" + err.Error())
	}
}
