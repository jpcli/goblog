package router

import (
	"goblog/utils/errors"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

func AppStart(addr string) {
	// gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	r.Static("/static/", "./static/")
	r.HTMLRender = loadTemplates()

	mgtAPI(r)
	mgtWeb(r)
	dispWeb(r)
	userWeb(r)

	err := r.Run(addr)
	panic(errors.WrapfErrorWithStack(err, "app has stoped"))
}

func loadTemplates() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	mgtView(&r)
	dispView(&r)
	userView(&r)
	return r
}
