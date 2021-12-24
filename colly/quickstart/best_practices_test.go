package quickstart

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"testing"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/colly/proxy"
)

// 如果任务足够复杂或具有不同类型的子任务,建议使用多个 Collector 来执行一个抓取作业
// Tips: 用Collector.ID调试以区分不同的 Collector
func TestMultiCollector(t *testing.T) {
	c1 := colly.NewCollector(
		colly.UserAgent("xy"),
		colly.AllowedDomains("foo.com", "bar.com"),
	)
	// Clone()复制具有相同配置但没有附加回调的 Collector
	c2 := c1.Clone()
	fmt.Printf("c1:%d, c2:%d", c1.ID, c2.ID)
}

// 在 Collector 之间传递自定义数据
func TestTransportData(t *testing.T) {
	c1 := colly.NewCollector()
	c2 := c1.Clone()

	// 使用 Collector 的Request()功能可以与其他 Collector 共享上下文
	c1.OnResponse(func(r *colly.Response) {
		r.Ctx.Put("v1", r.Headers.Get("Custom-Header"))
		c2.Request("GET", "https://foo.com/", nil, r.Ctx, nil)
	})
}

// 调试
// Colly 具有用于 Collector 调试的内置功能。提供了调试器接口和其他类型的调试器实现
// 你可以通过实现 debug.Debugger 接口实现自定义调试器, 如：debug.LogDebugger
func TestDebug(t *testing.T) {
	// 附加基本的日志调试器到Collector
	c := colly.NewCollector(
		colly.Debugger(&debug.LogDebugger{}),
	)
	fmt.Println("c", c)
}

var proxies []*url.URL = []*url.URL{
	&url.URL{Host: "127.0.0.1:8080"},
	&url.URL{Host: "127.0.0.1:8081"},
}

func randomProxySwitcher(_ *http.Request) (*url.URL, error) {
	return proxies[rand.Intn(len(proxies))], nil
}

// 分布式抓取可以根据抓取任务的要求以不同的方式实现
// 大多数时候,足以扩展网络通信层,这可以使用代理和 Colly 的代理切换器轻松实现
func TestProxy(t *testing.T) {
	/*
		当 HTTP 请求在多个代理之间分发时,使用代理切换器的抓取仍然保持集中。
		Colly 支持通过其SetProxyFunc()成员进行代理切换。
		可以SetProxyFunc()使用的签名将任何自定义函数传递给func(*http.Request) (*url.URL, error)。
		Tips:带有-D标志的 SSH 服务器可以用作 socks5 代理。
	*/
	// Colly具有内置的代理切换器,可根据每个请求轮流代理列表
	c := colly.NewCollector()
	if p, err := proxy.RoundRobinProxySwitcher(
		"socks5://127.0.0.1:1337",
		"socks5://127.0.0.1:1338",
		"http://127.0.0.1:8080",
	); err == nil {
		c.SetProxyFunc(p)
	}

	// 实现自定义代理切换器
	c.SetProxyFunc(randomProxySwitcher)
}

// Colly 能够用实现了storage.Storage 接口的任何存储后端替换默认的内存存储。检查现有存储。
// Colly 有一个内存中的存储后端,用于存储 cookie 和访问的 URL ,
// 但是任何实现 storage.Storage 的自定义存储后端都可以覆盖它
func TestStorage(t *testing.T) {
	/*
		现有的存储后端
			1.内存后端: Colly 的默认后端。使用Collector.SetStorage（）覆盖。
			2.Redis 后端:有关详细信息,请参见redis示例。
			3.boltdb 后端
			4.SQLite3 后端
			5.MongoDB 后端
			6.PostgreSQL 后端
	*/
}

// 扩展是 Colly 随附的小型帮助程序实用程序
// 以下示例启用了随机 User-Agent 切换器和 Referrer setter 扩展
// 并访问了 httpbin.org 两次。
func TestExtensions(t *testing.T) {
	c := colly.NewCollector()
	visited := false

	extensions.RandomUserAgent(c)
	extensions.Referer(c)

	c.OnResponse(func(r *colly.Response) {
		log.Println(string(r.Body))
		if !visited {
			visited = true
			r.Request.Visit("/get?q=2")
		}
	})

	c.Visit("http://httpbin.org/get")
}
