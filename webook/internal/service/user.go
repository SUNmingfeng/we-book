package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var ErrDuplicateEmail = repository.RepoErrDuplicateEmail

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
