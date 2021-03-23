package handler

import (
	"fmt"
	"goblog/service"
	"goblog/utils/markdown"

	"github.com/gin-gonic/gin"
)

func DispPost(c *gin.Context) {
	postID, err := getNumFromRegexp(`^(\d+)\.html$`, c.Param("path"))
	if postID < 1 || err != nil {
		HTMLNotFound(c, err)
		return
	}

	s := service.NewService()

	post, e := s.GetDispPostByPid(postID)
	if e != nil {
		handleServiceError(c, e)
		return
	}

	post.Text = markdown.ToHTML(post.Text)

	HTMLOK(c, "disp-post.html", gin.H{
		"post": post,
	}, &tkd{
		Title:       fmt.Sprintf("%s - 追风寻逸", post.Title),
		Description: post.Excerpt,
		Keywords:    post.Keywords,
	})

	// TODO: 在redis缓存中进行，使用incr，事务太耗费时间了
	go func() {
		// TODO: 错误处理
		_ = s.IncrPostViewCountByPid(postID)
	}()
}
