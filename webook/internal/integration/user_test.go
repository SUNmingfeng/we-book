package integration

import (
	"basic-go/webook/internal/integration/startup"
	"basic-go/webook/internal/web"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func TestUserService_SendSMSCode(t *testing.T) {
	rdb := startup.InitRedis()
	server := startup.InitWebserver()
	testcases := []struct {
		name string
		//提前准备数据
		before func(t *testing.T)
		//验证并删除数据
		after func(t *testing.T)
		phone string

		wantCode   int
		wantResult web.Result
	}{
		{
			name:   "发送成功的用例",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
				defer cancel()
				key := "phone_code:login:13299999999"
				code, err := rdb.Get(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, len(code) > 0)
				dur, err := rdb.TTL(ctx, key).Result()
				assert.NoError(t, err)
				assert.True(t, dur > time.Minute*9)
				err = rdb.Del(ctx, key).Err()
				assert.NoError(t, err)
			},
			phone:    "13299999999",
			wantCode: http.StatusOK,
			wantResult: web.Result{
				Msg: "发送成功",
			},
		},
		{
			name:     "手机号码为空",
			before:   func(t *testing.T) {},
			after:    func(t *testing.T) {},
			phone:    "",
			wantCode: http.StatusOK,
			wantResult: web.Result{
				Code: 4,
				Msg:  "手机号码不能为空",
			},
		},
		{
			name: "发送太频繁",
			before: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
				defer cancel()
				key := "phone_code:login:13299999999"
				err := rdb.Set(ctx, key, "123456", time.Minute*10).Err()
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second*1000)
				defer cancel()
				key := "phone_code:login:13299999999"
				code, err := rdb.GetDel(ctx, key).Result()
				assert.NoError(t, err)
				assert.Equal(t, "123456", code)
			},
			phone:    "13299999999",
			wantCode: http.StatusOK,
			wantResult: web.Result{
				Code: 4,
				Msg:  "短信发送频繁，请稍后再试",
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			tc.before(t)
			defer tc.after(t)
			req, err := http.NewRequest(http.MethodPost,
				"/users/login_sms/code/send",
				bytes.NewReader([]byte(fmt.Sprintf(`{"phone": "%s"}`, tc.phone))))
			req.Header.Set("Content-Type", "application/json")
			assert.NoError(t, err)
			recorder := httptest.NewRecorder()

			//执行
			server.ServeHTTP(recorder, req)

			assert.Equal(t, tc.wantCode, recorder.Code)
			if tc.wantCode != http.StatusOK {
				return
			}
			var result web.Result
			err = json.NewDecoder(recorder.Body).Decode(&result)
			assert.NoError(t, err)
			assert.Equal(t, tc.wantResult, result)
		})
	}
}
