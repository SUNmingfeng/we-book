package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	svcmocks "basic-go/webook/internal/service/mocks"
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserService_Signup(t *testing.T) {
	testcases := []struct {
		name string
		mock func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
		//预期的输入
		reqBuilder func(t *testing.T) *http.Request
		//预期的输出
		wantCode int
		wantBody string
	}{
		{
			name: "注册成功",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "test@test.com",
					PassWord: "test1234#",
				}).Return(nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"test@test.com",
				"passWord":"test1234#",
				"confirmPassWord":"test1234#"
			}`)))
				req.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "注册成功！",
		},
		{
			name: "bind出错",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"test@test.com",
				"passWord":"test1234#
			}`)))
				req.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "邮箱格式错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"test@",
				"passWord":"test1234#",
				"confirmPassWord":"test1234#"
			}`)))
				req.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "非法邮箱格式",
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"test@test.com",
				"passWord":"test1234#",
				"confirmPassWord":"test"
			}`)))
				req.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "两次密码输入不一致",
		},
		{
			name: "两次密码不一致",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				codeSvc := svcmocks.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"test@test.com",
				"passWord":"12345678",
				"confirmPassWord":"12345678"
			}`)))
				req.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "非法密码格式，至少包含字母、数字、特殊字符，1-9位",
		},
		{
			name: "系统错误",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "test@test.com",
					PassWord: "test1234#",
				}).Return(errors.New("db错误"))
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"test@test.com",
				"passWord":"test1234#",
				"confirmPassWord":"test1234#"
			}`)))
				req.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "系统错误",
		},
		{
			name: "注册邮箱冲突",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmocks.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domain.User{
					Email:    "test@test.com",
					PassWord: "test1234#",
				}).Return(ErrDuplicateEmail)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/users/signup", bytes.NewReader([]byte(`{
				"email":"test@test.com",
				"passWord":"test1234#",
				"confirmPassWord":"test1234#"
			}`)))
				req.Header.Add("Content-Type", "application/json")
				assert.NoError(t, err)
				return req
			},
			wantCode: http.StatusOK,
			wantBody: "注册邮箱冲突，请换一个",
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			//利用mock构造handler
			hdl := NewUserHandler(userSvc, codeSvc)
			//注册路由
			server := gin.Default()
			hdl.RegisterRoutes(server)

			//准备请求和recorder
			req := tc.reqBuilder(t)
			recorder := httptest.NewRecorder()
			//执行请求，响应存入recorder
			server.ServeHTTP(recorder, req)

			//判断
			assert.Equal(t, tc.wantCode, recorder.Code)
			assert.Equal(t, tc.wantBody, recorder.Body.String())
		})
	}
}
