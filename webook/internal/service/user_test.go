package service

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository"
	repmocks "basic-go/webook/internal/repository/mocks"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestPasswordEncrypt(t *testing.T) {
	password := []byte("1234hello#")
	encrypted, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	assert.NoError(t, err)
	println(string(encrypted))
	err = bcrypt.CompareHashAndPassword(encrypted, password)
	assert.NoError(t, err)
}

func Test_userService_Login(t *testing.T) {
	testcases := []struct {
		name string

		mock func(ctrl *gomock.Controller) repository.UserRepository

		//预期输入
		ctx      *gin.Context
		email    string
		password string
		//预期输出
		wantUser domain.User
		wantErr  error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRep := repmocks.NewMockUserRepository(ctrl)
				userRep.EXPECT().
					FindByEmail(gomock.Any(), "123456@qq.com").Return(domain.User{
					Email:    "123456@qq.com",
					PassWord: "$2a$10$oFaB1gq9awF0RZJbC/wn4eAXDLBpLHo3cK/nE5rwlYm9KgQ9mP2tG",
					Phone:    "1301234567",
				}, nil)
				return userRep
			},
			email:    "123456@qq.com",
			password: "1234hello#",
			wantUser: domain.User{
				Email:    "123456@qq.com",
				PassWord: "$2a$10$oFaB1gq9awF0RZJbC/wn4eAXDLBpLHo3cK/nE5rwlYm9KgQ9mP2tG",
				Phone:    "1301234567",
			},
			wantErr: nil,
		},
		{
			name: "未找到数据",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRep := repmocks.NewMockUserRepository(ctrl)
				userRep.EXPECT().
					FindByEmail(gomock.Any(), "123456@qq.com").Return(domain.User{}, repository.ErrRecordNotFound)
				return userRep
			},
			email:    "123456@qq.com",
			password: "1234hello#",
			wantErr:  ErrInvaildUserOrPassword,
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRep := repmocks.NewMockUserRepository(ctrl)
				userRep.EXPECT().
					FindByEmail(gomock.Any(), "123456@qq.com").Return(domain.User{}, errors.New("db错误"))
				return userRep
			},
			email:    "123456@qq.com",
			password: "1234hello#",
			wantErr:  errors.New("db错误"),
		},
		{
			name: "密码错误",
			mock: func(ctrl *gomock.Controller) repository.UserRepository {
				userRep := repmocks.NewMockUserRepository(ctrl)
				userRep.EXPECT().
					FindByEmail(gomock.Any(), "123456@qq.com").Return(domain.User{
					Email:    "123456@qq.com",
					PassWord: "$2a$10$oFaB1gq9awF0RZJbC/wn4eAXDLBpLHo3cK/nE5rwlYm9KgQ9mP2tG",
					Phone:    "1301234567",
				}, nil)
				return userRep
			},
			email:    "123456@qq.com",
			password: "1234hello#666",
			wantErr:  ErrInvaildUserOrPassword,
		},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRep := tt.mock(ctrl)
			userSvc := NewUserService(userRep)
			user, err := userSvc.Login(tt.ctx, tt.email, tt.password)
			assert.Equal(t, tt.wantUser, user)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
