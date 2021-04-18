package handler

import (
	"fmt"
	"goblog/service"
	"goblog/utils/option"
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
		Title:       fmt.Sprintf("网站声明 - %s", option.GetWebsiteName()),
		Description: fmt.Sprintf("%s网站声明", option.GetWebsiteName()),
		Keywords:    strings.Join([]string{"网站声明", option.GetWebsiteName()}, ","),
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
		Title:       fmt.Sprintf("文章归档 - %s", option.GetWebsiteName()),
		Description: fmt.Sprintf("%s所有文章归档", option.GetWebsiteName()),
		Keywords:    strings.Join([]string{"文章归档", option.GetWebsiteName()}, ","),
	})
}
