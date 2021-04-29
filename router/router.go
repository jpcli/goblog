package router

import (
	"github.com/gin-gonic/gin"
)

func AppRun() {
	r := gin.Default()

	// API路由
	api := r.Group("/api")
	apiRouter(api)

	r.Run()
}
