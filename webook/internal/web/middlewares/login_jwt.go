package middlewares

import (
	"basic-go/webook/internal/web"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
	"time"
)

type MiddlewareJWTBuilder struct {
}

func (m *MiddlewareJWTBuilder) CheckLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		//这两个接口不需要登录校验
		if path == "/users/signup" ||
			path == "/users/login" ||
			path == "/login_sms/code/send" ||
			path == "/login_sms" {
			println("跳过登录校验...")
			return
		}
		authCode := ctx.GetHeader("Authorization")
		if authCode == "" {
			//authorization不存在
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(authCode, " ")
		if len(segs) != 2 {
			//authorization内容错误
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr := segs[1]
		uc := web.UserClaims{}
		//uc传给了token.claims
		token, err := jwt.ParseWithClaims(tokenStr, &uc, func(token *jwt.Token) (interface{}, error) {
			return web.JWTKey, nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		if uc.UserAgent != ctx.GetHeader("User-Agent") {
			//后期埋点告警
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		expireTime := uc.ExpiresAt
		if expireTime.Before(time.Now()) {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if expireTime.Sub(time.Now()) < time.Second*50 {
			uc.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString(web.JWTKey)
			fmt.Println("现tokenstr：")
			fmt.Println(tokenStr)
			ctx.Header("x-jwt-token", tokenStr)
			if err != nil {
				fmt.Println(err)
			}
		}
		ctx.Set("user", uc)
	}
}
