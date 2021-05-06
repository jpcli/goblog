package middleware

import (
	"fmt"
	"goblog/utils/config"
	"goblog/utils/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 鉴权中间件
func Auth(c *gin.Context) {
	cookie, err := c.Request.Cookie("jwt")
	if err != nil {
		unauth(c, "请先登录")
		return
	}

	claims, err := jwt.ParseJWT(cookie.Value)
	if err != nil {
		unauth(c, "登录凭证无效或过期")
		return
	}

	if claims.User != config.OauthGithubAdminUser() {
		unauth(c, "")
		return
	}

	c.Next()
}

// 未授权
func unauth(c *gin.Context, msg string) {
	res := gin.H{
		"code": -1,
	}
	if msg == "" {
		res["msg"] = "您没有权限访问该页面"
	} else {
		res["msg"] = fmt.Sprintf("您没有权限访问该页面，%s", msg)
	}

	c.JSON(http.StatusUnauthorized, res)
	c.Abort()
}
