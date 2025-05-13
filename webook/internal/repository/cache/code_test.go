package cache

import (
	"basic-go/webook/internal/repository/cache/redismocks"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

func TestRedisCodeCache_Set(t *testing.T) {
	keyFunc := func(biz, phone string) string {
		return fmt.Sprintf("phone_code:%s:%s", biz, phone)
	}
	testcases := []struct {
		name  string
		mock  func(ctrl *gomock.Controller) redis.Cmdable
		ctx   context.Context
		biz   string
		phone string
		code  string
		want  error
	}{
		{
			name: "set成功",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				ud := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				ud.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("验证码", "15511111111")}, []any{"222222"}).Return(cmd)
				return ud
			},
			ctx:   context.Background(),
			biz:   "验证码",
			phone: "15511111111",
			code:  "222222",
			want:  nil,
		},
		{
			name: "redis报错",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				ud := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis error"))
				ud.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("验证码", "15511111111")}, []any{"222222"}).Return(cmd)
				return ud
			},
			ctx:   context.Background(),
			biz:   "验证码",
			phone: "15511111111",
			code:  "222222",
			want:  errors.New("redis error"),
		},
		{
			name: "验证码不存在过期时间",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				ud := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-2))
				ud.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("验证码", "15511111111")}, []any{"222222"}).Return(cmd)
				return ud
			},
			ctx:   context.Background(),
			biz:   "验证码",
			phone: "15511111111",
			code:  "222222",
			want:  ErrKeyNotExist,
		},
		{
			name: "发送太频繁",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				ud := redismocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-1))
				ud.EXPECT().Eval(gomock.Any(), luaSetCode, []string{keyFunc("验证码", "15511111111")}, []any{"222222"}).Return(cmd)
				return ud
			},
			ctx:   context.Background(),
			biz:   "验证码",
			phone: "15511111111",
			code:  "222222",
			want:  ErrCodeSendTooMany,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			c := NewCodeCache(tc.mock(ctrl))
			err := c.Set(tc.ctx, tc.biz, tc.phone, tc.code)
			assert.Equal(t, tc.want, err)
		})
	}
}
