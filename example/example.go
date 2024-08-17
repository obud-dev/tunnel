package main

import (
	"github.com/gin-gonic/gin"
	"github.com/obud-dev/tunnel/pkg/utils"
)

func main() {
	utils.InitLogger()
	go utils.PrintMemoryUsage()

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "world",
		})
	})
	r.Run(":8080")
}
