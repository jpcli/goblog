package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// api错误返回
func apiFailed(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  msg,
	})
	c.Abort()
}

// api错误输入
func apiErrorInput(c *gin.Context) {
	apiFailed(c, "请求的参数有误")
}

// api成功，msg可选，但只取第一个
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

// api未授权
func apiUnauth(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"code": -1,
		"msg":  "您没有权限访问该页面",
	})
	c.Abort()
}
