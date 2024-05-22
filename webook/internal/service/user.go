package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail        = repository.RepoErrDuplicateEmail
	ErrInvaildUserOrPassword = errors.New("用户不存在或密码不正确")
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}

}

func (svc *UserService) Signup(ctx *gin.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx *gin.Context, email string, password string) (domain.User, error) {
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

func (svc *UserService) FindById(ctx *gin.Context, userid int64) (domain.User, error) {
	return svc.repo.FindById(ctx, userid)
}

func (svc *UserService) UpdateInfo(ctx *gin.Context, user domain.User) error {
	return svc.repo.UpdateFields(ctx, user)
}
