package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/dao"
	"context"
	"github.com/gin-gonic/gin"
	"time"
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
		Birthday: u.Birthday.UnixMilli(),
	})
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx *gin.Context, userid int64) (domain.User, error) {
	u, err := repo.dao.FindById(ctx, userid)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
		PassWord: u.PassWord,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
	}
}

func (repo *UserRepository) UpdateFields(ctx *gin.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *UserRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		ID:       user.Id,
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,
	}
}
