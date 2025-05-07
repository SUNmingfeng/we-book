package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.RepoErrDuplicateUser
	ErrInvaildUserOrPassword = errors.New("用户不存在或密码不正确")
)

type UserService interface {
	Signup(ctx *gin.Context, u domain.User) error
	Login(ctx *gin.Context, email string, password string) (domain.User, error)
	FindById(ctx *gin.Context, userid int64) (domain.User, error)
	UpdateInfo(ctx *gin.Context, user domain.User) error
	FindOrCreate(ctx *gin.Context, phone string) (domain.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Signup(ctx *gin.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *userService) Login(ctx *gin.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrRecordNotFound {
		return u, ErrInvaildUserOrPassword
	}
	if err != nil {
		return u, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.PassWord), []byte(password))
	if err == repository.ErrRecordNotFound {
		return u, ErrInvaildUserOrPassword
	}
	if err != nil {
		return u, err
	}
	return u, nil
}

func (svc *userService) FindById(ctx *gin.Context, userid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userid)
}

func (svc *userService) UpdateInfo(ctx *gin.Context, user domain.User) error {
	return svc.repo.UpdateFields(ctx, user)
}

func (svc *userService) FindOrCreate(ctx *gin.Context, phone string) (domain.User, error) {
	u, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrRecordNotFound {
		//两种情况：一种是err为其他错误，统一为系统错误，另一种是找到数据，err==nil
		return u, err
	}
	err = svc.repo.Create(ctx, domain.User{
		Phone: phone,
	})
	if err != nil && err != repository.RepoErrDuplicateUser {
		return domain.User{}, err
	}
	//err == nil (创建用户成功)或 err == RepoErrDuplicateUser(用户存在)
	//这里查找用户不一定能找到，因为数据库可能存在主从延迟
	return svc.repo.FindByPhone(ctx, phone)
}
