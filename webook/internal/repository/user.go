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

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	UpdateFields(ctx *gin.Context, user domain.User) error
	FindByEmail(ctx *gin.Context, email string) (domain.User, error)
	FindById(ctx context.Context, userid int64) (domain.User, error)
	FindByPhone(ctx *gin.Context, phone string) (domain.User, error)
	FindByWechat(ctx *gin.Context, openid string) (domain.User, error)
}
type CachedUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewCachedUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CachedUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CachedUserRepository) Create(ctx context.Context, u domain.User) error {
	//domain的user不一定和dao的user完全对应，所以dao中要重新定义一个自己的user，需要的数据从domain的user中传入
	return repo.dao.Insert(ctx, repo.toEntity(u))
}

func (repo *CachedUserRepository) FindByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) FindById(ctx context.Context, userid int64) (domain.User, error) {
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
		//go func() {
		//	err = repo.cache.Set(ctx, du)
		//	if err != nil {
		//		// 网络崩溃或者redis崩溃的情况下，缓存中依然没有存入该条数据
		//		log.Println(err)
		//	}
		//}()
		err = repo.cache.Set(ctx, du)
		if err != nil {
			// 网络崩溃或者redis崩溃的情况下，缓存中依然没有存入该条数据
			log.Println(err)
		}
		return du, nil
	default:
		// 访问redis，降级，防止缓存穿透，不再查数据库，返回空
		return domain.User{}, err
	}
}

func (repo *CachedUserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Nickname: u.Nickname,
		Email:    u.Email.String,
		Phone:    u.Phone.String,
		PassWord: u.PassWord,
		Birthday: time.UnixMilli(u.Birthday),
		AboutMe:  u.AboutMe,
		Ctime:    time.UnixMilli(u.Ctime),
		WechatInfo: domain.WechatInfo{
			OpenId:  u.WechatOpenId.String,
			UnionId: u.WechatUnionId.String,
		},
	}
}

func (repo *CachedUserRepository) UpdateFields(ctx *gin.Context, user domain.User) error {
	return repo.dao.UpdateById(ctx, repo.toEntity(user))
}

func (repo *CachedUserRepository) toEntity(user domain.User) dao.User {
	return dao.User{
		Id: user.Id,
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
		WechatOpenId: sql.NullString{
			String: user.WechatInfo.OpenId,
			Valid:  user.WechatInfo.OpenId != "",
		},
		WechatUnionId: sql.NullString{
			String: user.WechatInfo.UnionId,
			Valid:  user.WechatInfo.UnionId != "",
		},
	}
}

func (repo *CachedUserRepository) FindByPhone(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := repo.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *CachedUserRepository) FindByWechat(ctx *gin.Context, openId string) (domain.User, error) {
	u, err := repo.dao.FindByWechat(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}
