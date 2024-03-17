package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	"github.com/gin-gonic/gin"
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
	err := svc.repo.Create(ctx, u)
	if err != nil {
		return err
	}
	return nil
}
