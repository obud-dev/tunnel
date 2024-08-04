package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/obud-dev/tunnel/pkg/svc"
)

func ApiServer(ctx *svc.ServerCtx) {
	r := gin.Default()
	r.Use(AuthMiddleware)

	r.GET("/tunnels", func(c *gin.Context) {
		// tunnels := ctx.Db.GetTunnels()
		// c.JSON(http.StatusOK, tunnels)
	})
	r.Run(ctx.Config.Api)
}

func AuthMiddleware(c *gin.Context) {
	// 从请求中获取 basic auth
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}
	// 验证 basic auth
	if username != "admin" || password != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized",
		})
		c.Abort()
		return
	}
	c.Next()
}
