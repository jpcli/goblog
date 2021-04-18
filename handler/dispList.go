package handler

import (
	"fmt"
	"goblog/repository"
	"goblog/service"
	"goblog/utils/option"
	"strings"

	"github.com/gin-gonic/gin"
)

func DispIndexPostList(c *gin.Context) {
	page, err := getNumFromRegexp(`^(\d+)$`, c.Param("path"))
	if page < 1 || err != nil {
		HTMLNotFound(c, err)
		return
	}

	s := service.NewService()

	postList, pageNav, e := s.GetDispPostList(page, option.GetEachPageLimit())
	if e != nil {
		handleServiceError(c, e)
		return
	}

	pageNav.URLFormat = "/p/%d/"

	// TODO: tkd设置在option里
	var title string
	if page == 1 {
		title = fmt.Sprintf("%s - 分享生活，创造智慧！", option.GetWebsiteName())
	} else {
		title = fmt.Sprintf("第%d页 - 文章列表 - %s", page, option.GetWebsiteName())
	}

	HTMLOK(c, "disp-post-list.html", gin.H{
		"postList":    postList,
		"pageNavInfo": pageNav,
	}, &tkd{
		Title:       title,
		Description: fmt.Sprintf("%s是一个分享生活与技术、创造智慧的个人笔记博客", option.GetWebsiteName()),
		Keywords:    strings.Join([]string{"网页前端", "Web后端", "Python", "技术分享", option.GetWebsiteName()}, ","),
	})
}

func DispCategoryPostList(c *gin.Context) {
	page, err := getNumFromRegexp(`^(\d+)$`, c.Param("path"))
	if page < 1 || err != nil {
		HTMLNotFound(c, err)
		return
	}

	s := service.NewService()

	slug := c.Param("slug")
	postList, pageNav, term, e := s.GetDispPostListBySlug(repository.TermTypeCategory, slug, page, option.GetEachPageLimit())
	if e != nil {
		handleServiceError(c, e)
		return
	}

	pageNav.URLFormat = fmt.Sprintf("/category/%s/%%d/", slug)

	var title string
	if page == 1 {
		title = fmt.Sprintf("%s - 分类 - %s", term.Name, option.GetWebsiteName())
	} else {
		title = fmt.Sprintf("%s - 第%d页 - 分类 - %s", term.Name, page, option.GetWebsiteName())
	}

	HTMLOK(c, "disp-post-list.html", gin.H{
		"postList":    postList,
		"pageNavInfo": pageNav,
	}, &tkd{
		Title:       title,
		Description: term.Description,
		Keywords:    strings.Join([]string{term.Name, option.GetWebsiteName()}, ","),
	})
}

func DispTagPostList(c *gin.Context) {
	page, err := getNumFromRegexp(`^(\d+)$`, c.Param("path"))
	if page < 1 || err != nil {
		HTMLNotFound(c, err)
		return
	}

	s := service.NewService()

	slug := c.Param("slug")
	postList, pageNav, term, e := s.GetDispPostListBySlug(repository.TermTypeTag, slug, page, option.GetEachPageLimit())
	if e != nil {
		handleServiceError(c, e)
		return
	}

	pageNav.URLFormat = fmt.Sprintf("/tag/%s/%%d/", slug)

	var title string
	if page == 1 {
		title = fmt.Sprintf("%s - 标签 - %s", term.Name, option.GetWebsiteName())
	} else {
		title = fmt.Sprintf("%s - 第%d页 - 标签 - %s", term.Name, page, option.GetWebsiteName())
	}

	HTMLOK(c, "disp-post-list.html", gin.H{
		"postList":    postList,
		"pageNavInfo": pageNav,
	}, &tkd{
		Title:       title,
		Description: term.Description,
		Keywords:    strings.Join([]string{term.Name, option.GetWebsiteName()}, ","),
	})
}
