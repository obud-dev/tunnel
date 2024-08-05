package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/obud-dev/tunnel/pkg/model"
	"github.com/obud-dev/tunnel/pkg/response"
	"github.com/obud-dev/tunnel/pkg/svc"
)

func ApiServer(ctx *svc.ServerCtx) {
	r := gin.Default()
	r.Use(AuthMiddleware)

	r.GET("/tunnels", func(c *gin.Context) {
		tunnels, err := ctx.TunnelModel.GetTunnels()
		response.Response(c, tunnels, err)
	})

	r.POST("/tunnels", func(c *gin.Context) {
		var tunnel model.Tunnel
		if err := c.BindJSON(&tunnel); err != nil {
			response.Response(c, nil, response.New(-1, err.Error()))
			return
		}
		err := ctx.TunnelModel.Insert(&tunnel)
		response.Response(c, nil, err)
	})

	r.GET("/tunnels/:id", func(c *gin.Context) {
		id := c.Param("id")
		tunnel, err := ctx.TunnelModel.GetTunnelByID(id)
		response.Response(c, tunnel, err)
	})

	r.PUT("/tunnels/:id", func(c *gin.Context) {
		id := c.Param("id")
		var tunnel model.Tunnel
		if err := c.BindJSON(&tunnel); err != nil {
			response.Response(c, nil, response.New(-1, err.Error()))
			return
		}
		tunnel.ID = id
		err := ctx.TunnelModel.Update(&tunnel)
		response.Response(c, nil, err)
	})

	r.DELETE("/tunnels/:id", func(c *gin.Context) {
		id := c.Param("id")
		tunnel, err := ctx.TunnelModel.GetTunnelByID(id)
		if err != nil {
			response.Response(c, nil, err)
			return
		}
		err = ctx.TunnelModel.Delete(tunnel)
		response.Response(c, nil, err)
	})

	fmt.Printf("API Server is running on %s\n\n", ctx.Config.Api)
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
