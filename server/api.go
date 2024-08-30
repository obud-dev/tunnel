package main

import (
	"embed"
	"net/http"
	"time"

	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/obud-dev/tunnel/pkg/config"
	"github.com/obud-dev/tunnel/pkg/model"
	"github.com/obud-dev/tunnel/pkg/response"
	"github.com/obud-dev/tunnel/pkg/svc"
	"github.com/obud-dev/tunnel/pkg/utils"
	"github.com/rs/zerolog/log"
)

//go:embed web/dist/*
var staticFiles embed.FS

//go:embed web/dist/index.html
var indexHtml []byte

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
		tunnel.ID = utils.GenerateID()
		tunnel.Token = utils.GenerateID()[0:32]
		tunnel.Uptime = time.Now().Unix()
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
		tunnel.Uptime = time.Now().Unix()
		err := ctx.TunnelModel.Update(&tunnel)
		response.Response(c, nil, err)
	})

	api.DELETE("/tunnels/:id", func(c *gin.Context) {
		id := c.Param("id")
		var routes []model.Route

		tunnel, err := ctx.TunnelModel.GetTunnelByID(id)
		if err != nil {
			response.Response(c, nil, err)
			return
		}
		routes, err = ctx.RouteModel.GetRoutesByTunnelID(id)
		if err != nil {
			response.Response(c, nil, err)
			return
		}
		for _, route := range routes {
			if route.TunnelID == id {
				err = ctx.RouteModel.Delete(&route)
				if err != nil {
					response.Response(c, nil, err)
					return
				}
			}
		}
		ctx.DelTunnel(tunnel.ID)
		err = ctx.TunnelModel.Delete(tunnel)
		response.Response(c, nil, err)
	})

	api.POST("/tunnels/:id/refreshtoken", func(c *gin.Context) {
		id := c.Param("id")
		tunnel, err := ctx.TunnelModel.GetTunnelByID(id)
		if err != nil {
			response.Response(c, nil, err)
			return
		}
		// 32位随机字符串
		token := utils.GenerateID()[0:32]
		tunnel.Token = token
		ctx.DelTunnel(tunnel.ID)
		err = ctx.TunnelModel.Update(tunnel)
		response.Response(c, token, err)
	})

	api.GET("/routes/:tid", func(c *gin.Context) {
		tid := c.Param("tid")
		routes, err := ctx.RouteModel.GetRoutesByTunnelID(tid)
		response.Response(c, routes, err)
	})

	api.POST("/routes", func(c *gin.Context) {
		var route model.Route
		if err := c.BindJSON(&route); err != nil {
			response.Response(c, nil, response.New(-1, err.Error()))
			return
		}
		route.ID = utils.GenerateID()
		err := ctx.RouteModel.Insert(&route)
		ctx.UpdateRoutes()
		response.Response(c, nil, err)
	})

	api.DELETE("/routes/:id", func(c *gin.Context) {
		id := c.Param("id")
		route, err := ctx.RouteModel.GetRouteByID(id)
		if err != nil {
			response.Response(c, nil, err)
			return
		}
		err = ctx.RouteModel.Delete(route)
		ctx.UpdateRoutes()
		response.Response(c, nil, err)
	})

	api.PUT("/routes", func(c *gin.Context) {
		var route model.Route
		if err := c.BindJSON(&route); err != nil {
			response.Response(c, nil, response.New(-1, err.Error()))
			return
		}
		err := ctx.RouteModel.Update(&route)
		ctx.UpdateRoutes()
		response.Response(c, nil, err)
	})

	api.GET("/token/:tid", func(c *gin.Context) {
		tid := c.Param("tid")
		tunnel, err := ctx.TunnelModel.GetTunnelByID(tid)
		if err != nil {
			response.Response(c, nil, err)
			return
		}
		config := &config.ClientConfig{
			TunnelID: tunnel.ID,
			Token:    tunnel.Token,
			Server:   ctx.Config.Host + ctx.Config.ListenOn,
		}
		token, err := config.Encode()
		if err != nil {
			response.Response(c, nil, err)
			return
		}
		response.Response(c, token, err)
	})

	api.GET("/server/info", func(c *gin.Context) {
		server, err := ctx.ServerModel.GetServer()
		response.Response(c, server, err)
	})

	r.Use(static.Serve("/", static.EmbedFolder(staticFiles, "web/dist")))
	r.NoRoute(func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", indexHtml)
	})

	log.Info().Msgf("API Server is running on %s", ctx.Config.Api)
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
