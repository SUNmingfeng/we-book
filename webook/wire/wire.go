//go:build wireinject

package wire

import (
	"basic-go/webook/wire/repository"
	"basic-go/webook/wire/repository/dao"
	"github.com/google/wire"
)

func InitUserRepository() *repository.UserRepository {
	wire.Build(repository.NewUserRepository, dao.NewUserDAO, InitDB)
	return &repository.UserRepository{}
}
