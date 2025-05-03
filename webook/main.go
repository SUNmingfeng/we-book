package main

import (
	"basic-go/webook/config"
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/sms"
	"basic-go/webook/internal/service/sms/localsms"
	"basic-go/webook/internal/web"
	"basic-go/webook/internal/web/middlewares"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	sess_redis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"strings"
	"time"
)

func main() {
	db := initDB()
	server := initWebServer()
	redisClient := redis.NewClient(&redis.Options{
		Addr: config.Config.Redis.Addr,
	})

	codeSvc := initCodeService(redisClient)
	initUserHdl(db, server, redisClient, codeSvc)
	//server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello! 启动成功！")
	})
	//server.GET("/user/:name", func(ctx *gin.Context) {
	//	name := ctx.Param("name")
	//	ctx.String(http.StatusOK, name)
	//})
	//
	//server.GET("/order", func(ctx *gin.Context) {
	//	id := ctx.Query("id")
	//	ctx.String(http.StatusOK, "我的ID是："+id)
	//})
	server.Run(":8080") // 默认也是8080
}

func initUserHdl(db *gorm.DB, server *gin.Engine, redisClient redis.Cmdable, codeSvc *service.CodeService) {
	ud := dao.NewUserDAO(db)
	uc := cache.NewUserCache(redisClient)
	ur := repository.NewUserRepository(ud, uc)
	us := service.NewUserService(ur)
	hdl := web.NewUserHandler(us, codeSvc)

	hdl.RegisterRoutes(server)
}

func initCodeService(redisClient redis.Cmdable) *service.CodeService {
	cc := cache.NewCodeCache(redisClient)
	crepo := repository.NewCodeRepository(cc)
	return service.NewCodeService(crepo, initMemorySms())
}

func initMemorySms() sms.Service {
	return localsms.NewService()
}
func initWebServer() *gin.Engine {
	server := gin.Default()
	// 在Use中注册的方法都是middleware（use中的参数为HandlerFunc，以*context为参数的方法都可以作为HandlerFunc），进入engine的请求都会先执行middleware
	server.Use(
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
		},
	)

	//redisClient := redis.NewClient(&redis.Options{
	//	Addr: config.Config.Redis.Addr,
	//})
	//server.Use(ratelimit.NewBuilder(redisClient, time.Second, 100).Build())
	//useSession(server)
	useJWT(server)
	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DNS))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
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

func useJWT(server *gin.Engine) {
	login := &middlewares.MiddlewareJWTBuilder{}
	server.Use(login.CheckLogin())
}
