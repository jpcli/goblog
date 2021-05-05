package router

import (
	"goblog/utils/config"
	"path"

	"github.com/gin-gonic/gin"
)

// 运行应用程序，执行前应已经正常打开所有依赖组件
func AppRun() {
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

	r.Run(config.AppAddr())
}
