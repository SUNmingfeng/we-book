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
	// 在Use中注册的方法都是middleware（use中的参数为HandlerFunc，以*context为参数的方法都可以作为HandlerFunc），进入engine的请求都会先执行middleware
	server.Use(
		//跨域问题是由于发请求的协议+域名+端口和接收请求的协议+域名+端口对不上
		//配置跨域策略
		cors.New(cors.Config{
			//允许带的认证信息，cookie等
			AllowCredentials: true,
			//允许header中带的头
			AllowHeaders: []string{"Content-Type", "authorization"},
			//origin：请求来源
			AllowOriginFunc: func(origin string) bool {
				//允许本地
				if strings.Contains(origin, "localhost") {
					return true
				}
				//允许公司域名
				return strings.Contains(origin, "the_company_com")
			},
			//允许的访问方法类型，不配时是默认全部允许
			//AllowMethods: []string{"POST"},
			//preflight检测后的有效期
			MaxAge: 12 * time.Hour,
		}),
		//以*context为参数的方法都可以作为HandlerFunc
		func(ctx *gin.Context) {
			println("这里是middleware")
		})
	hdl.RegisterRoutes(server)
	server.Run(":8080")
}
