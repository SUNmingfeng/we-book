package repository

import (
	"basic-go/webook/wire/repository/dao"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(ud *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: ud,
	}
}
