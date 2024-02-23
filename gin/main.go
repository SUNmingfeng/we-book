package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	server := gin.Default()
	server.GET("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello,world")
	})
	server.POST("/login", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello login")
	})

	server.POST("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "hello "+name)
	})
	server.GET("/order", func(ctx *gin.Context) {
		id := ctx.Query("id")
		ctx.String(http.StatusOK, "order ID is :"+id)
	})
	server.Run(":8080")

}
