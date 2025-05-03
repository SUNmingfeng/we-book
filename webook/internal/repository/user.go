package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	"basic-go/webook/internal/repository/dao"
	"context"
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
	"time"
)

var (
	RepoErrDuplicateUser = dao.ErrDuplicateEmail
	ErrRecordNotFound    = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao   *dao.UserDAO
	cache *cache.UserCache
}

func NewUserRepository(dao *dao.UserDAO, cache *cache.UserCache) *UserRepository {
	return &UserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	//domain的user不一定和dao的user完全对应，所以dao中要重新定义一个自己的user，需要的数据从domain的user中传入
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *UserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) FindById(ctx context.Context, userid int64) (domain.User, error) {
	du, err := repo.cache.Get(ctx, userid)
	switch err {
	case nil:
		return du, nil
	case cache.ErrorKeyNotExist:
		// redis正常，没查到key
		u, err := repo.dao.FindById(ctx, userid)
		if err != nil {
			return domain.User{}, err
		}
		du = repo.toDomain(u)
		// 回写缓存
		go func() {
			err = repo.cache.Set(ctx, du)
			if err != nil {
				// 网络崩溃或者redis崩溃的情况下，缓存中依然没有存入该条数据
				log.Println(err)
			}
		}()
		return du, nil
	default:
		// 访问redis，降级，防止缓存穿透，不再查数据库，返回空
		return domain.User{}, err
	}
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
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
		ID: user.Id,
		Email: sql.NullString{
			String: user.Email,
			Valid:  user.Email != "",
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  user.Phone != "",
		},
		Nickname: user.Nickname,
		Birthday: user.Birthday.UnixMilli(),
		AboutMe:  user.AboutMe,
	}
}

func (repo *UserRepository) FindByPhone(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
