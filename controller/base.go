package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func apiFailed(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  msg,
	})
	c.Abort()
}

func apiErrorInput(c *gin.Context) {
	apiFailed(c, "请求的参数有误")
}

func apiOK(c *gin.Context, data gin.H, msg ...string) {
	m := ""
	if len(msg) > 0 {
		m = msg[0]
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  m,
		"data": data,
	})
	c.Abort()
}

func apiUnauth(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": -1,
		"msg":  "您没有权限访问该页面",
	})
	c.Abort()
}
