package web

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	time2 "time"
)

type jwtHandler struct {
}

func (j *jwtHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		Uid: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time2.Now().Add(time2.Minute * 30)),
		},
		UserAgent: ctx.GetHeader("User-Agent"),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenString, err := token.SignedString(JWTKey)
	if err != nil {
		ctx.String(http.StatusOK, "jwt token生成错误")
	}
	ctx.Header("x-jwt-token", tokenString)
}

var JWTKey = []byte("ehOk5JYoP2glSsMXmSvhRdupSr9ToEuiMLvSmKU127SpCkCxDB8JoMgONCszg55N")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
