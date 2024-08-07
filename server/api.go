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
	r.Use(AuthMiddleware(ctx))

	api := r.Group("/api")
	api.GET("/tunnels", func(c *gin.Context) {
		tunnels, err := ctx.TunnelModel.GetTunnels()
		response.Response(c, tunnels, err)
	})

	api.POST("/tunnels", func(c *gin.Context) {
		var tunnel model.Tunnel
		if err := c.BindJSON(&tunnel); err != nil {
			response.Response(c, nil, response.New(-1, err.Error()))
			return
		}
		err := ctx.TunnelModel.Insert(&tunnel)
		response.Response(c, nil, err)
	})

	api.GET("/tunnels/:id", func(c *gin.Context) {
		id := c.Param("id")
		tunnel, err := ctx.TunnelModel.GetTunnelByID(id)
		response.Response(c, tunnel, err)
	})

	api.PUT("/tunnels/:id", func(c *gin.Context) {
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

	api.DELETE("/tunnels/:id", func(c *gin.Context) {
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

func AuthMiddleware(ctx *svc.ServerCtx) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从请求中获取 basic auth
		username, password, ok := c.Request.BasicAuth()
		if !ok {
			c.Header("WWW-Authenticate", "Basic realm=Authorization Required")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		// 验证 basic auth
		if username != ctx.Config.User || password != ctx.Config.Password {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
