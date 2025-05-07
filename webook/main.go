package main

import (
	"basic-go/webook/internal/web/middlewares"
	"github.com/gin-contrib/sessions"
	sess_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := InitWebserver()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello! 启动成功！")
	})
	server.Run(":8080") // 默认也是8080
}

func useSession(server *gin.Engine) {
	login := &middlewares.MiddlewareBuilder{}
	//用来存储session
	//store := cookie.NewStore([]byte("ehOk5JYoP2glSsMXmSvhRdupSr9TgEuiMLvSmKU127SpCkCxDB8JoMgONCszg55N"))
	//store := memstore.NewStore([]byte("ehOk5JYoP2glSsMXmSvhRdupSr9TgEuiMLvSmKU127SpCkCxDB8JoMgONCszg55N"), []byte("ehOk5JYoP2glSsMXmSvhRdupSr9TgEuiMLvSmKU127SpCkCxDB8JoMgONCszg55y"))
	store, err := sess_redis.NewStore(16, "tcp", "localhost:6379", "",
		[]byte("yMvUN8X2MdHYBoF8Dvi60SMjXCe4aD9k"),
		[]byte("yMvUN8X2MdHYBoF8Dvi60SMjXCe4aD9b"))
	if err != nil {
		panic(err)
	}
	//先初始化session，后面才能使用
	server.Use(sessions.Sessions("ssid", store), login.CheckLogin())
}
