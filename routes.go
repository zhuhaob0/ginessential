package main

import (
	"net/http"
	"oceanlearn/ginessential/controller"

	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})

	r.POST("api/auth/register", controller.Register)

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "请求路径或方法错误",
		})
	})

	return r
}
