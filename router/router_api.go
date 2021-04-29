package router

import (
	"goblog/controller"

	"github.com/gin-gonic/gin"
)

func apiRouter(r *gin.RouterGroup) {
	r.POST("/post/detail", controller.PostAdd) // 新增文章
	r.GET("/post/detail/:id")                  // 获取文章
	r.PUT("/post/detail/:id")                  // 修改文章
	r.PATCH("/post/status/:id")                // 修改文章状态
	r.DELETE("/post/:id")                      // 删除文章
	r.GET("/post/list")                        // 获取文章列表，参数?page=&limit=

	r.POST("/term/detail")    // 新增项
	r.GET("/term/detail/:id") // 获取项
	r.PUT("/term/detail/:id") // 修改项
	r.GET("/term/list")       // 获取项列表，参数?page=&limit=

	r.POST("/upload") // 上传附件

	r.PUT("/option/base") // 修改基本配置（全修改）
	r.GET("/option/base") // 获取基本配置
}
