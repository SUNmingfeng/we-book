package web

import (
	"basic-go/webook/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserHandler struct {
	svc *service.UserService
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
	if err := ctx.Bind(&req); err != nil {
		return
	}
	ctx.String(http.StatusOK, "hello,little wang")
}

func (h *UserHandler) Login(ctx *gin.Context) {

}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) ProFile(ctx *gin.Context) {

}
