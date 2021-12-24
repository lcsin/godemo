package example

import (
	"fmt"
	"log"
	"testing"

	"github.com/gocolly/colly"
	"github.com/gocolly/colly/queue"
	"github.com/gocolly/redisstorage"
)

func TestRedisBackend(t *testing.T) {

	urls := []string{
		"http://httpbin.org/",
		"http://httpbin.org/ip",
		"http://httpbin.org/cookies/set?a=b&c=d",
		"http://httpbin.org/cookies",
	}

	c := colly.NewCollector()

	// create the redis storage
	storage := &redisstorage.Storage{
		Address:  "10.225.137.237:6379",
		Password: "root",
		DB:       0,
		Prefix:   "httpbin_test",
	}

	// add storage to the collector
	err := c.SetStorage(storage)
	if err != nil {
		panic(err)
	}

	// delete previous data from storage
	if err := storage.Clear(); err != nil {
		log.Fatal(err)
	}

	// close redis client
	defer storage.Client.Close()

	// create a new request queue with redis storage backend
	q, _ := queue.New(2, storage)

	c.OnResponse(func(r *colly.Response) {
		log.Println("Cookies:", c.Cookies(r.Request.URL.String()))
	})
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	// add URLs to the queue
	for _, u := range urls {
		q.AddURL(u)
	}
	// consume requests
	q.Run(c)

}
