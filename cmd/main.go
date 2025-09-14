package main

import (
	"github.com/YuanJey/go-log/pkg/log"
	"github.com/gin-gonic/gin"
	"zf-server/internal/api/game"
	"zf-server/internal/api/middleware"
	"zf-server/internal/api/order"
	"zf-server/internal/api/statics"
	"zf-server/internal/api/user"
	"zf-server/pkg/common/config"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery())
	r.LoadHTMLGlob("./static/*")
	r.Static("/tools", "./tools")
	staticGroup := r.Group("/push")
	{
		staticGroup.GET("/register", statics.Register)
		staticGroup.GET("/push_order", statics.PushOrder)
		staticGroup.GET("/login", statics.Login)
		staticGroup.GET("/index", statics.Index)
	}
	pullGroup := r.Group("/pull")
	{
		pullGroup.GET("/pull_order", statics.PullOrder)
	}
	// 用户相关路由（无需认证）
	userGroup := r.Group("/user")
	{
		userGroup.POST("/register", user.Register)
		userGroup.POST("/login", user.Login) // 添加登录路由
	}

	gameGroup := r.Group("/game")
	{
		gameGroup.POST("/pull", game.PullOrder)
	}

	// 需要认证的路由组
	authGroup := r.Group("/")
	authGroup.Use(middleware.AuthMiddleware()) // 应用认证中间件
	{
		authUserGroup := authGroup.Group("/auth")
		{
			authUserGroup.GET("/info", user.UserInfo)
		}
		balanceGroup := authGroup.Group("/admin")
		{
			balanceGroup.POST("/recharge", user.Recharge)
		}
		ordergroup := authGroup.Group("order")
		{
			ordergroup.POST("/push", order.PushOrder)
		}
	}

	Port := "8000"
	if config.Config.Server.Port != "" {
		Port = config.Config.Server.Port
	}
	log.Debug("server port is " + Port)
	address := "0.0.0.0:" + Port
	err := r.Run(address)
	if err != nil {
		panic("api start failed " + err.Error())
	}
}
