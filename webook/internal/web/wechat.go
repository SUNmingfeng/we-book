package web

import (
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/oauth2/wechat"
	"github.com/gin-gonic/gin"
	"net/http"
)

type OAuth2WechatHandler struct {
	jwtHandler
	svc  wechat.Service
	user service.UserService
}

func NewOAuth2WechatHandler(svc wechat.Service, user service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:  svc,
		user: user,
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.auth2URL)
	g.Any("callback", o.callback)
}

func (o *OAuth2WechatHandler) auth2URL(ctx *gin.Context) {
	val, err := o.svc.AuthURL(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "生成跳转URL失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})
}

func (o *OAuth2WechatHandler) callback(ctx *gin.Context) {
	code := ctx.Query("code")
	//state := ctx.Query("state")
	dw, err := o.svc.VerifyCode(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "校验授权码失败",
		})
		return
	}
	du, err := o.user.FindOrCreateByWechat(ctx, dw)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
	}
	o.setJWTToken(ctx, du.Id)
	ctx.JSON(http.StatusOK, Res{
		Msg: "登录成功",
	})

}
