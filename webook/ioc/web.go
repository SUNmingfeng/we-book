package ioc

import (
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middlewares"
	"basic-go/webook/pkg/ginx/middleware/ratelimit"
	"basic-go/webook/pkg/limiter"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler, OAuth2WechatHandler *web.OAuth2WechatHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	OAuth2WechatHandler.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(redisClient redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		//跨域问题是由于发请求的协议+域名+端口和接收请求的协议+域名+端口对不上
		//配置跨域策略
		cors.New(cors.Config{
			//允许携带认证信息，cookie等
			AllowCredentials: true,
			//允许header中带的头
			AllowHeaders:  []string{"Content-Type", "Authorization"},
			ExposeHeaders: []string{"x-jwt-token"},
			//origin：请求来源
			AllowOriginFunc: func(origin string) bool {
				//允许本地
				if strings.HasPrefix(origin, "http://localhost") {
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
		},
		//ratelimit.NewBuilder(ratelimit.NewRedisSlidingWindowLimiter(redisClient, time.Second, 1000)).Build(),
		ratelimit.NewBuilder(limiter.NewRedisSlidingWindowsLimiter(redisClient, time.Second, 1000)).Build(),
		(&middlewares.LoginJWTMiddlewareBuilder{}).CheckLogin(),
	}
}
