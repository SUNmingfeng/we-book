package main

import (
	"basic-go/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func main() {
	hdl := web.NewUserHandler()
	server := gin.Default()
	// 在Use中注册的方法都是middleware，进入engine的请求都会先执行middleware
	server.Use(cors.New(cors.Config{ //临时解决跨域问题
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "authorization"},
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return true
		},
		MaxAge: 12 * time.Hour,
	}), func(ctx *gin.Context) {
		println("这里是middleware")
	})
	hdl.RegisterRoutes(server)
	server.Run(":8080")
}
