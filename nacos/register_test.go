package nacos

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// 客户端和服务都配置
// 可以配置多个ServerConfig，客户端回对这些服务端做轮询请求
var (
	cc = constant.ClientConfig{
		NamespaceId:         "2ac20d2d-aa9f-49f2-9365-1022bff1ada9",
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "tmp\\nacos\\log",
		CacheDir:            "tmp\\nacos\\cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "debug",
	}
	sc = []constant.ServerConfig{
		{
			IpAddr:      "192.168.5.130",
			ContextPath: "/nacos",
			Port:        8848,
			Scheme:      "http",
		},
	}

	err error
	// 创建服务发现客户端
	ns naming_client.INamingClient
	// 创建动态配置客户端
	cs config_client.IConfigClient
)

func init() {
	// 初始化服务发现客户端
	ns, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic("[create naming client failed] " + err.Error())
	}

	// 初始化动态配置客户端
	cs, err = clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic("[crate config client failed]" + err.Error())
	}
}

// 注册实例
func TestRegisterInstance(t *testing.T) {
	instance, err := ns.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          "10.225.137.237",
		Port:        8080,
		ServiceName: "demo.go",
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{"idc": "shanghai"},
		ClusterName: "cluster-a", // 默认值DEFAULT
		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic("[registry instance failed] " + err.Error())
	}
	fmt.Println(instance)
}

// 注销实例
func TestDeregisterInstance(t *testing.T) {
	success, err := ns.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          "10.225.137.237",
		Port:        8080,
		ServiceName: "demo.go",
		Ephemeral:   true,
		Cluster:     "cluster-a", // 默认值DEFAULT
		GroupName:   "group-a",   // 默认值DEFAULT_GROUP
	})
	if err != nil {
		panic("[deregister instance failed]" + err.Error())
	}
	fmt.Println(success)
}

// 获取服务信息
func TestGetService(t *testing.T) {
	service, err := ns.GetService(vo.GetServiceParam{
		Clusters:    []string{"mxshop-srvs"},
		ServiceName: "user-srv",
		GroupName:   "mxshop",
	})
	if err != nil {
		panic("[get service info failed] " + err.Error())
	}
	fmt.Println(service)
}

// 获取所有的实例列表
func TestSelectAllInstance(t *testing.T) {
	allInstances, err := ns.SelectAllInstances(vo.SelectAllInstancesParam{
		Clusters:    []string{"cluster-a"},
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})
	if err != nil {
		panic("[get all instance failed]" + err.Error())
	}
	fmt.Println(allInstances)
}

// 获取实例列表
func TestSelectInstances(t *testing.T) {
	instances, err := ns.SelectInstances(vo.SelectInstancesParam{
		Clusters:    []string{"cluster-a"},
		ServiceName: "demo.go",
		GroupName:   "group-a",
		HealthyOnly: true,
	})
	if err != nil {
		panic("[get instances failed]" + err.Error())
	}
	fmt.Println(instances)
}

// 获取一个健康的实例
func TestSelectOneHealthInstance(t *testing.T) {
	healthyInstance, err := ns.SelectOneHealthyInstance(vo.SelectOneHealthInstanceParam{
		Clusters:    []string{"cluster-a"},
		ServiceName: "demo.go",
		GroupName:   "group-a",
	})
	if err != nil {
		panic("[get healthy instance failed]" + err.Error())
	}
	fmt.Println(healthyInstance)
}

// 监听服务变化
func TestSubscribe(t *testing.T) {
	err := ns.Subscribe(&vo.SubscribeParam{
		ServiceName: "user-srv",
		GroupName:   "mxshop",                // 默认值DEFAULT_GROUP
		Clusters:    []string{"mxshop-srvs"}, // 默认值DEFAULT
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			fmt.Printf("\n callback return services:%v \n\n", services)
			//fmt.Println(services[0].Ip, services[0].Port, services[0].ServiceName)
		},
	})
	if err != nil {
		panic("[subscribe instance change failed] " + err.Error())
	}
	time.Sleep(time.Second * 300)
}

// 取消服务监听
func TestUnsubscribe(t *testing.T) {
	err := ns.Unsubscribe(&vo.SubscribeParam{
		ServiceName: "demo.go",
		GroupName:   "group-a",             // 默认值DEFAULT_GROUP
		Clusters:    []string{"cluster-a"}, // 默认值DEFAULT
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			log.Printf("\n\n callback return services:%v \n\n", services)
		},
	})
	if err != nil {
		panic("[unsubscribe instance change failed] " + err.Error())

	}
}

// 获取所有实例信息
func TestGetAllServicesInfo(t *testing.T) {
	servicesInfo, err := ns.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: "",
		GroupName: "group-a",
		PageNo:    1,
		PageSize:  10,
	})
	if err != nil {
		panic("[get all service info failed]" + err.Error())
	}
	fmt.Println(servicesInfo)
}
