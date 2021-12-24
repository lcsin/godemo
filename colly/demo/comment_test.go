package demo

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"testing"
)

func TestBilibiliComment(t *testing.T) {
	url := "https://www.bilibili.com/video/BV1h64y1h7rb"

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		fmt.Println(string(r.Body))
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Fatalln(err)
	})

	err := c.Visit(url)
	if err != nil {
		log.Fatalln(err)
	}
}
