package router

import (
	"fmt"
	"goblog/handler"
	"goblog/utils/config"
	"goblog/utils/errors"
	"path/filepath"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// mgtAPI sets all routers for management system api.
func mgtAPI(r *gin.Engine) {
	api := r.Group(config.AppMgtURI()+"/api", handler.MgtAccessControl)

	api.POST("/getPost", handler.GetPostAPI)
	api.POST("/getPostList", handler.GetPostListAPI)
	api.POST("/editPost", handler.EditPostAPI)
	api.POST("/modifyPostStatus", handler.ModifyPostStatus)

	api.POST("/getTerm", handler.GetTermAPI)
	api.POST("/getTermList", handler.GetTermListAPI)
	api.POST("/editTerm", handler.EditTermAPI)

	api.POST("/upload", handler.UploadAPI)
}

func mgtWeb(r *gin.Engine) {
	web := r.Group(config.AppMgtURI()+"/view", handler.MgtAccessControl)

	web.GET("/*path", handler.MgtViewHandler)
}

// management system view
func mgtView(r *multitemplate.Renderer) {
	(*r).AddFromFiles("mgt-index.html", "./view/mgt/index.html")

	includes, err := filepath.Glob("./view/mgt/includes/*.html")
	if err != nil {
		panic(errors.WrapfErrorWithStack(err, "failed to load management system view files"))
	}

	for _, include := range includes {
		files := []string{}
		files = append(files, "./view/mgt/frame.html", include)
		(*r).AddFromFiles(fmt.Sprintf("mgt-%s", filepath.Base(include)), files...)
	}
}
