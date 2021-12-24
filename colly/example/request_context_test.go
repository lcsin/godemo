package example

import (
	"fmt"
	"testing"

	"github.com/gocolly/colly"
)

func TestRequestContext(t *testing.T) {
	// Instantiate default collector
	c := colly.NewCollector()

	// Before making a request put the URL with
	// the key of "url" into the context of the request
	c.OnRequest(func(r *colly.Request) {
		r.Ctx.Put("url", r.URL.String()+"lcx")
	})

	// After making a request get "url" from
	// the context of the request
	c.OnResponse(func(r *colly.Response) {
		fmt.Println(r.Ctx.Get("url"))
	})

	// Start scraping on https://en.wikipedia.org
	c.Visit("https://httpbin.org/delay/2")
}
