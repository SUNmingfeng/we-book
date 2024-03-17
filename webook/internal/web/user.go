package web

import (
	"basic-go/webook/internal/domain"
	"basic-go/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexpPattern = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$"
	//至少包含字母、数字、特殊字符，1-9位
	passwordRegexpPattern = "^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[^\\da-zA-Z\\s]).{1,9}$"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordPexExp *regexp.Regexp
	svc            *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		//预编译正则
		emailRexExp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		passwordPexExp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// 使用group分组路由来简化注册
	ug := server.Group("/users")
	ug.POST("/signup", h.SginUp)
	ug.POST("/login", h.Login)
	ug.GET("/profile", h.ProFile)
	ug.POST("/edit", h.Edit)
}

func (h *UserHandler) SginUp(ctx *gin.Context) {
	type SignReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignReq
	//取出前端数据到req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	//对前端输出的一些格式验证
	isEmail, _ := h.emailRexExp.MatchString(req.Email)
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}
	isPassword, _ := h.passwordPexExp.MatchString(req.Password)
	if !isPassword {
		ctx.String(http.StatusOK, "非法密码格式，至少包含字母、数字、特殊字符，1-9位")
		return
	}

	if req.ConfirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次密码输入不一致")
		return
	}

	//执行登陆，校验用户和密码

	//存储数据
	err := h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		PassWord: req.Password,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	ctx.String(http.StatusOK, fmt.Sprint("注册成功！"))
}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) ProFile(ctx *gin.Context) {

}
