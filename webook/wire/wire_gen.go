// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"basic-go/webook/wire/repository"
	"basic-go/webook/wire/repository/dao"
)

// Injectors from wire.go:

func InitUserRepository() *repository.UserRepository {
	db := InitDB()
	userDAO := dao.NewUserDAO(db)
	userRepository := repository.NewUserRepository(userDAO)
	return userRepository
}
