package web

import (
	"basic-go/webook/internal/service"
	"basic-go/webook/internal/service/oauth2/wechat"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/lithammer/shortuuid/v4"
	"net/http"
)

type OAuth2WechatHandler struct {
	jwtHandler
	svc           wechat.Service
	user          service.UserService
	key           []byte
	cookStateName string
}

func NewOAuth2WechatHandler(svc wechat.Service, user service.UserService) *OAuth2WechatHandler {
	return &OAuth2WechatHandler{
		svc:           svc,
		user:          user,
		key:           []byte("ehOk5JYoP2glSsMXmSvhRdupSr9ToEuiMLvSmKU127SpCkCxDB8JoMgONCszg55Q"),
		cookStateName: "jwt-state",
	}
}

func (o *OAuth2WechatHandler) RegisterRoutes(server *gin.Engine) {
	g := server.Group("/oauth2/wechat")
	g.GET("/authurl", o.auth2URL)
	g.Any("Callback", o.Callback)
}

func (o *OAuth2WechatHandler) auth2URL(ctx *gin.Context) {
	state := uuid.New()
	val, err := o.svc.AuthURL(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "生成跳转URL失败",
		})
		return
	}
	err = o.setStateCookie(ctx, state)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "服务器异常",
		})
	}
	ctx.JSON(http.StatusOK, Result{
		Data: val,
	})
}

func (o *OAuth2WechatHandler) Callback(ctx *gin.Context) {
	err := o.verifyState(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "非法请求",
		})
	}
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

func (o *OAuth2WechatHandler) verifyState(ctx *gin.Context) error {
	state := ctx.Query("state")
	ck, err := ctx.Cookie(o.cookStateName)
	if err != nil {
		return fmt.Errorf("cookie获取state、失败: %v", err)
	}
	var uc StateClaims
	_, err = jwt.ParseWithClaims(ck, &uc, func(token *jwt.Token) (interface{}, error) {
		return o.key, nil
	})
	if err != nil {
		return fmt.Errorf("解析jwt token失败: %v", err)
	}
	if state != uc.State {
		return fmt.Errorf("state校验不一致")
	}
}

func (o *OAuth2WechatHandler) setStateCookie(ctx *gin.Context, state string) error {
	sc := StateClaims{
		State: state,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, sc)
	tokenString, err := token.SignedString(o.key)
	if err != nil {
		return err
	}
	ctx.SetCookie(o.cookStateName, tokenString, 600, "/oauth2/wechat/Callback", "", false, true)
	return nil
}

type StateClaims struct {
	jwt.RegisteredClaims
	State string `json:"state"`
}
