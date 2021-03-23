package handler

import (
	"fmt"
	"goblog/utils/config"
	"goblog/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

//TODO：recover的中间件：记录日志，输出信息（栈信息等），例如在APIError后不return，直接panic错误，用中间件接受

func APIError(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  msg,
		"data": "",
	})
	c.Abort()
}

func APIOK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "",
		"data": data,
	})
	c.Abort()
}

func MgtAccessControl(c *gin.Context) {
	e := func(c *gin.Context, msg string) {
		c.JSON(http.StatusForbidden, gin.H{
			"code": -1,
			"msg":  fmt.Sprintf("error: %s", msg),
			"data": "",
		})
		c.Abort()
	}

	cookie, err := c.Request.Cookie("jwt")
	if err != nil {
		e(c, "please login first")
		return
	}

	claims, err := jwt.ParseJWT(cookie.Value)
	if err != nil {
		e(c, "token is invalid or expired")
		return
	}

	if claims.GithubUser != config.OauthAdmin() {
		e(c, "has no access to view this page")
		return
	}

	// if claims.Version < 1 {
	// 	panic(fmt.Errorf("jwt is expired"))
	// }

	c.Next()
}
