//go:build wireinject

package startup

import (
	"basic-go/webook/internal/repository"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/web"
	"basic-go/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

func InitWebserver() *gin.Engine {
	wire.Build(
		ioc.InitDB, InitRedis,
		dao.NewUserDAO,
		cache.NewUserCache, cache.NewCodeCache,
		repository.NewCachedUserRepository, repository.NewCachedCodeRepository,
		service.NewUserService, service.NewCodeService,
		ioc.InitSMSService,
		web.NewUserHandler,
		ioc.InitGinMiddlewares,
		ioc.InitWebServer,
	)
	return gin.Default()
}
