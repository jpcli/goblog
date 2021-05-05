package router

import (
	"goblog/utils/config"
	"path"

	"github.com/gin-gonic/gin"
)

func AppRun(addr string) {
	r := gin.Default()

	// 静态资源
	r.Static("/static", "./static")

	// 管理页路由
	admin := r.Group(path.Join("/admin", config.AdminSafetyFactor()))
	{
		// API路由
		api := admin.Group("/api")
		apiRouter(api)
		// 视图路由
		admin.Static("/view", "./view/admin")
	}

	r.Run(addr)
}
