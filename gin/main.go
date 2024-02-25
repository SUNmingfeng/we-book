package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//创建web服务器
	server := gin.Default()
	//路由注册
	//context：处理请求，返回响应
	//静态路由
	server.GET("/hello", func(ctx *gin.Context) {
		//把响应写回前端
		ctx.String(http.StatusOK, "hello,world")
	})
	server.POST("/login", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello login")
	})

	//参数路由
	server.POST("/users/:name", func(ctx *gin.Context) {
		name := ctx.Param("name")
		ctx.String(http.StatusOK, "hello "+name)
	})

	//查询参数
	server.GET("/order", func(ctx *gin.Context) {
		id := ctx.Query("id")
		ctx.String(http.StatusOK, "order ID is :"+id)
	})

	//通配符路由
	server.GET("/view/*aaa", func(ctx *gin.Context) {
		path := ctx.Param("aaa") //
		ctx.String(http.StatusOK, "html is :"+path)
	})
	server.Run(":8080") //默认是8080

}
