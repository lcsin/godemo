package example

import (
	"log"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
)

func TestLogin(t *testing.T) {
	// create a new collector
	c := colly.NewCollector()

	// authenticate
	err := c.Post("http://localhost:8080/login", map[string]string{"username": "admin", "password": "admin"})
	if err != nil {
		log.Fatal("login failed:", err.Error())
	}

	// attach callbacks after login
	c.OnResponse(func(r *colly.Response) {
		log.Println(string(r.Body))
		log.Println("response received", r.StatusCode)
	})

	// start scraping
	c.Visit("http://localhost:8080")
}

func TestLoginServer(t *testing.T) {
	r := gin.Default()
	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		if username == "admin" && password == "admin" {
			c.JSON(http.StatusOK, gin.H{
				"code": http.StatusOK,
				"msg":  "login success",
			})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code": http.StatusNotFound,
			"msg":  "login failed, account not found",
		})
	})

	r.Run(":8080")
}
