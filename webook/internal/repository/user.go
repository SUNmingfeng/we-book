package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/dao"
	"context"
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	//domain的user不一定和dao的user完全对应，所以dao中要重新定义一个自己的user，需要的数据从domain的user中传入
	err := repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		PassWord: u.PassWord,
	})
	if err != nil {
		return err
	}
	return nil
}
