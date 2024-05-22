package middlewares

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MiddlewareBuilder struct {
}

func (m *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		//这两个接口不需要登录校验
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		sess := sessions.Default(context)
		if sess.Get("userId") == nil {
			//中断，不再向后执行
			context.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
