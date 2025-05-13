package repository

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/repository/cache"
	cachemocks "basic-go/webook/internal/repository/cache/mocks"
	"basic-go/webook/internal/repository/dao"
	daomocks "basic-go/webook/internal/repository/dao/mocks"
	"context"
	"database/sql"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func TestCachedUserRepository_FindById(t *testing.T) {
	testCases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache)

		ctx    context.Context
		userId int64

		wantUser domain.User
		wantErr  error
	}{
		{
			name: "查找成功，缓存未命中",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDao := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				uid := int64(123)
				userCache.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrorKeyNotExist)
				userDao.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{
					Id:       uid,
					Nickname: "test1",
					Email: sql.NullString{
						String: "test@example.com",
						Valid:  true,
					},
					PassWord: "12345",
					Birthday: 123,
					AboutMe:  "自我介绍",
					Phone: sql.NullString{
						String: "19900000000",
						Valid:  true,
					},
					Ctime: 101,
					Utime: 102,
				}, nil)
				userCache.EXPECT().Set(gomock.Any(), domain.User{
					Id:       123,
					Email:    "test@example.com",
					Nickname: "test1",
					PassWord: "12345",
					Birthday: time.UnixMilli(123),
					AboutMe:  "自我介绍",
					Phone:    "19900000000",
					Ctime:    time.UnixMilli(101),
				}).Return(nil)
				return userDao, userCache
			},
			ctx:    context.Background(),
			userId: 123,
			wantUser: domain.User{
				Id:       123,
				Email:    "test@example.com",
				Nickname: "test1",
				PassWord: "12345",
				Birthday: time.UnixMilli(123),
				AboutMe:  "自我介绍",
				Phone:    "19900000000",
				Ctime:    time.UnixMilli(101),
			},
			wantErr: nil,
		},
		{
			name: "未找到用户",
			mock: func(ctrl *gomock.Controller) (dao.UserDAO, cache.UserCache) {
				userDao := daomocks.NewMockUserDAO(ctrl)
				userCache := cachemocks.NewMockUserCache(ctrl)
				uid := int64(123)
				userCache.EXPECT().Get(gomock.Any(), uid).Return(domain.User{}, cache.ErrorKeyNotExist)
				userDao.EXPECT().FindById(gomock.Any(), uid).Return(dao.User{}, dao.ErrRecordNotFound)
				return userDao, userCache
			},
			ctx:      context.Background(),
			userId:   123,
			wantUser: domain.User{},
			wantErr:  ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			ud, uc := tc.mock(ctrl)
			svc := NewCachedUserRepository(ud, uc)
			user, err := svc.FindById(tc.ctx, tc.userId)
			assert.Equal(t, tc.wantErr, err)
			assert.Equal(t, tc.wantUser, user)
		})
	}
}
