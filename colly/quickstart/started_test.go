package quickstart

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/gocolly/colly"
)

// Colly的收集器的创建及回调
func TestStarted(t *testing.T) {
	// Colly的主要实体对象，管理网络通信并负责指向附加的回调
	c := colly.NewCollector()

	// 主要的回调
	// Called before a request
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Called if error occured during the request
	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	// Called after response received
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})

	// Called right after OnResponse if the received content is HTML
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})
	c.OnHTML("tr td:nth-of-type(1)", func(e *colly.HTMLElement) {
		fmt.Println("First column of a table row:", e.Text)
	})

	// Called right after OnHTML if the received content is HTML or XML
	c.OnXML("//h1", func(e *colly.XMLElement) {
		fmt.Println(e.Text)
	})

	// Called after OnXML callbacks
	c.OnScraped(func(r *colly.Response) {
		fmt.Println("Finished", r.Request.URL)
	})
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString() string {
	b := make([]byte, rand.Intn(10)+10)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// Colly收集器的配置
func TestConfiguration(t *testing.T) {
	// 1.通过传入配置创建Collector
	_ = colly.NewCollector(
		colly.UserAgent("xy"),
		colly.AllowURLRevisit(),
	)

	// 2.通过覆盖Collector的属性，可以在抓取作业的任何时候更改配置
	c2 := colly.NewCollector()
	// 会在每个请求上更改User-Agent
	c2.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", RandomString())
	})

	/*
		3.通过环境变量更改Collector的默认配置，而无需重新编译
		  环境解析是Collector初始化的最后一步，初始化之后的每个配置更改
		  都会覆盖从环境解析的配置
			COLLY_ALLOWED_DOMAINS (comma separated list of domains)
			COLLY_CACHE_DIR (string)
			COLLY_DETECT_CHARSET (y/n)
			COLLY_DISABLE_COOKIES (y/n)
			COLLY_DISALLOWED_DOMAINS (comma separated list of domains)
			COLLY_IGNORE_ROBOTSTXT (y/n)
			COLLY_FOLLOW_REDIRECTS (y/n)
			COLLY_MAX_BODY_SIZE (int)
			COLLY_MAX_DEPTH (int - 0 means infinite)
			COLLY_PARSE_HTTP_ERROR_RESPONSE (y/n)
			COLLY_USER_AGENT (string)
	*/

	// 4.HTTP配置
	c3 := colly.NewCollector()
	// colly使用Golang的默认http客户端作为网络层。
	// 可以通过更改默认的HTTP round tripper来调整HTTP选项
	c3.WithTransport(&http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	})
}
