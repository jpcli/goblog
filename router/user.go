package router

import (
	"fmt"
	"goblog/handler"
	"goblog/utils/errors"
	"net/http"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func userWeb(r *gin.Engine) {
	oauth := r.Group("/oauth")
	oauth.POST("/github", handler.GithubOauth)

	r.GET("login.html", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user-login.html", gin.H{})
	})
}

func userView(r *multitemplate.Renderer) {
	// 非模板嵌套
	htmls, err := filepath.Glob("./view/user/*.html")
	if err != nil {
		panic(errors.WrapfErrorWithStack(err, "failed to load htmls in user system"))
	}
	for _, html := range htmls {
		(*r).AddFromGlob(fmt.Sprintf("user-%s", filepath.Base(html)), html)
	}
}
