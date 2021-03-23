package handler

import (
	"goblog/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Sitemap(c *gin.Context) {
	s := service.NewService()
	sitemap, e := s.GetSitemap()
	if e != nil {
		handleServiceError(c, e)
		return
	}

	c.Header("Content-Type", "application/xml")
	c.String(http.StatusOK, sitemap)
}

func Declaration(c *gin.Context) {
	HTMLOK(c, "disp-declaration.html", gin.H{}, &tkd{
		Title:       "网站声明 - 追风寻逸",
		Description: "追风寻逸网站声明",
		Keywords:    strings.Join([]string{"网站声明", "追风寻逸"}, ","),
	})
}

func Archives(c *gin.Context) {
	s := service.NewService()
	archives, e := s.GetArchives()
	if e != nil {
		handleServiceError(c, e)
		return
	}

	HTMLOK(c, "disp-archives.html", gin.H{
		"archives": archives,
	}, &tkd{
		Title:       "文章归档 - 追风寻逸",
		Description: "追风寻逸所有文章归档",
		Keywords:    strings.Join([]string{"文章归档", "追风寻逸"}, ","),
	})
}
