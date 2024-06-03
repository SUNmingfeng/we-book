package middlewares

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type MiddlewareBuilder struct {
}

func (m *MiddlewareBuilder) CheckLogin() gin.HandlerFunc {
	gob.Register(time.Now())
	return func(context *gin.Context) {
		path := context.Request.URL.Path
		//这两个接口不需要登录校验
		if path == "/users/signup" || path == "/users/login" {
			return
		}
		sess := sessions.Default(context)
		userId := sess.Get("userId")
		if userId == nil {
			//中断，不再向后执行
			context.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		const updateTime = "update_time"
		val := sess.Get(updateTime)
		lastUpdateTime, ok := val.(time.Time)
		now := time.Now()
		if val == nil || !ok || now.Sub(lastUpdateTime) > time.Second*10 {
			sess.Set(updateTime, now)
			sess.Set("userId", userId)
			err := sess.Save()
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
