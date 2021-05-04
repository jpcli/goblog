package router

import (
	"github.com/gin-gonic/gin"
)

func AppRun() {
	r := gin.Default()

	// 静态资源
	r.Static("/static", "./static")

	// 管理页路由
	admin := r.Group("/admin")
	{
		// API路由
		api := admin.Group("/api")
		apiRouter(api)
		// 视图路由
		admin.Static("/view", "./view/admin")
	}

	r.Run()
}
