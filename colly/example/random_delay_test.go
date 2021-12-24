package example

import (
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/debug"
	"testing"
	"time"
)

func TestRandomDelay(t *testing.T) {
	url := "https://httpbin.org/delay/2"

	// Instantiate default collector
	c := colly.NewCollector(
		// Attach a debugger to the collector
		colly.Debugger(&debug.LogDebugger{}),
		colly.Async(true),
	)

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println(err)
	})

	// Limit the number of threads started by colly to two
	// when visiting links which domains' matches "*httpbin.*" glob
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*httpbin.*",
		Parallelism: 2,
		RandomDelay: 5 * time.Second,
	})

	// Start scraping in four threads on https://httpbin.org/delay/2
	for i := 0; i < 4; i++ {
		c.Visit(fmt.Sprintf("%s?n=%d", url, i))
	}
	// Start scraping on https://httpbin.org/delay/2
	c.Visit(url)
	// Wait until threads are finished
	c.Wait()
}
