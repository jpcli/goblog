package router

import (
	"goblog/controller"

	"github.com/gin-gonic/gin"
)

func apiRouter(r *gin.RouterGroup) {
	r.POST("/post/detail", controller.PostAdd)               // 新增文章
	r.GET("/post/detail/:id", controller.PostGet)            // 获取文章
	r.PUT("/post/detail/:id", controller.PostModify)         // 修改文章
	r.PATCH("/post/status/:id", controller.PostStatusModify) // 修改文章状态
	r.DELETE("/post/:id")                                    // 删除文章（软删除）
	r.GET("/post/list/normal", controller.PostListNormal)    // 获取正常文章列表
	r.GET("/post/list/trash")                                // 获取回收站文章列表

	r.POST("/term/detail", controller.TermAdd)       // 新增项
	r.GET("/term/detail/:id", controller.TermGet)    // 获取项
	r.PUT("/term/detail/:id", controller.TermModify) // 修改项
	r.DELETE("/term/:id")                            // 删除项
	r.GET("/term/list", controller.TermList)         // 获取项列表

	r.POST("/upload", controller.Upload) // 上传附件

	r.PUT("/option/base", controller.BaseOptionSet) // 修改基本配置（全修改）
	r.GET("/option/base", controller.BaseOptionGet) // 获取基本配置
}
