package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/dao"
	"context"
	"github.com/gin-gonic/gin"
)

var (
	RepoErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrRecordNotFound     = dao.ErrRecordNotFound
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
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		PassWord: u.PassWord,
	})
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.ID,
		Email:    u.Email,
		PassWord: u.PassWord,
	}
}
