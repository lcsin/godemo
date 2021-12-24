package nacos

import (
	"fmt"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/vo"
)

// 发布配置
func TestPublishConfig(t *testing.T) {
	success, err := cs.PublishConfig(vo.ConfigParam{
		DataId:  "dataId",
		Group:   "group",
		Content: "hello world!222222"})
	if err != nil {
		panic("[publish config failed] " + err.Error())
	}
	fmt.Println(success)
}

// 删除配置
func TestDeleteConfig(t *testing.T) {
	success, err := cs.DeleteConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group"})
	if err != nil {
		panic("[delete config failed] " + err.Error())
	}
	fmt.Println(success)
}

// 获取配置
func TestGetConfig(t *testing.T) {
	content, err := cs.GetConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group"})
	if err != nil {
		panic("[get config failed] " + err.Error())
	}
	fmt.Println(content)
}

// 监听配置变化
func TestListenConfig(t *testing.T) {
	err := cs.ListenConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
		OnChange: func(namespace, group, dataId, data string) {
			fmt.Println("\ngroup:" + group + ", dataId:" + dataId + ", data:\n" + data)
		},
	})
	if err != nil {
		panic("[listen config failed] " + err.Error())
	}
	time.Sleep(time.Second * 300)
}

// 取消配置监听
func TestCancelListenConfig(t *testing.T) {
	err := cs.CancelListenConfig(vo.ConfigParam{
		DataId: "dataId",
		Group:  "group",
	})
	if err != nil {
		panic("[cancel listen config failed] " + err.Error())
	}
}

// 搜索配置
func TestSearchConfig(t *testing.T) {
	configPage, err := cs.SearchConfig(vo.SearchConfigParam{
		Search:   "blur",
		DataId:   "",
		Group:    "",
		PageNo:   1,
		PageSize: 10,
	})
	if err != nil {
		panic("[search config failed] " + err.Error())
	}
	fmt.Println(configPage)
}
