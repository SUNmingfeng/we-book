package web

import (
	"basic-go/webook/internal/service"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	emailRegexpPattern    = "^[a-zA-Z0-9_-]+@[a-zA-Z0-9_-]+(.[a-zA-Z0-9_-]+)+$"
	passwordRegexpPattern = "^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[^\\da-zA-Z\\s]).{1,9}$"
)

type UserHandler struct {
	svc            *service.UserService
	emailRexExp    *regexp.Regexp
	passwordPexExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		//预编译正则
		emailRexExp:    regexp.MustCompile(emailRegexpPattern, regexp.None),
		passwordPexExp: regexp.MustCompile(passwordRegexpPattern, regexp.None),
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
	//由bind完成回写
	if err := ctx.Bind(&req); err != nil {
		return
	}
	isEmail, _ := h.emailRexExp.MatchString(req.Email)
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}
	isPassword, _ := h.passwordPexExp.MatchString(req.Password)
	if !isPassword {
		ctx.String(http.StatusOK, "非法密码格式")
		return
	}
	ctx.String(http.StatusOK, fmt.Sprintf("你正在注册：%v", req))
}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) ProFile(ctx *gin.Context) {

}
