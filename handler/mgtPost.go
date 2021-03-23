package handler

import (
	"fmt"
	"goblog/repository"
	"goblog/service"
	"math"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/spf13/cast"
)

func GetPostAPI(c *gin.Context) {
	pid, err := cast.ToUint32E(c.DefaultQuery("id", "0"))
	if err != nil {
		APIError(c, fmt.Sprintf("unexpected id value: %s", err))
		return
	}

	s := service.NewService()

	post, err := s.GetMgtPostByPid(pid)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get post: %s", err))
		return
	}

	allCategories, err := s.GetMgtTermList(repository.TermTypeCategory, 1, math.MaxUint32)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get all categories: %s", err))
		return
	}

	allTags, err := s.GetMgtTermList(repository.TermTypeTag, 1, math.MaxUint32)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get all tags: %s", err))
		return
	}

	APIOK(c, gin.H{
		"post":          post,
		"allCategories": allCategories,
		"allTags":       allTags,
	})
}

func GetPostListAPI(c *gin.Context) {
	page, err := cast.ToUint32E(c.DefaultPostForm("page", "1"))
	if err != nil {
		APIError(c, fmt.Sprintf("unexpected page value: %s", err))
		return
	}

	limit, err := cast.ToUint32E(c.DefaultPostForm("limit", "10"))
	if err != nil {
		APIError(c, fmt.Sprintf("unexpected limit value: %s", err))
		return
	}

	s := service.NewService()

	postList, err := s.GetMgtPostList(page, limit)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get post list: %s", err))
		return
	}

	count, err := s.CountPost()
	if err != nil {
		APIError(c, fmt.Sprintf("failed to count post: %s", err))
		return
	}

	APIOK(c, gin.H{
		"count": count,
		"list":  postList,
	})
}

func EditPostAPI(c *gin.Context) {
	pid, err := cast.ToUint32E(c.DefaultQuery("id", "0"))
	if err != nil {
		APIError(c, fmt.Sprintf("unexpected id value: %s", err))
		return
	}

	var data service.MgtPost
	err = c.ShouldBindBodyWith(&data, binding.JSON)
	if err != nil {
		APIError(c, fmt.Sprintf("miss post data: %s", err))
		return
	}

	s := service.NewService()

	id, err := s.EditPost(pid, &data)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to edit term: %s", err))
		return
	}

	APIOK(c, gin.H{
		"id": id,
	})
}

func ModifyPostStatus(c *gin.Context) {
	pid, err := cast.ToUint32E(c.Query("id"))
	if err != nil {
		APIError(c, fmt.Sprintf("unexpected post id: %s", err))
		return
	} else if pid == 0 {
		APIError(c, "unexpected post id: 0")
		return
	}

	d := struct {
		Status string `json:"status" binding:"required"`
	}{}
	err = c.ShouldBindJSON(&d)
	if err != nil {
		APIError(c, "miss status data")
		return
	} else if st := d.Status; st != repository.PostStatusPublish && st != repository.PostStatusSticky && st != repository.PostStatusTrash {
		APIError(c, fmt.Sprintf("unexpected post status: %s", st))
		return
	}

	s := service.NewService()
	err = s.ModifyPostStatus(pid, d.Status)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to modify post status: %s", err))
		return
	}

	APIOK(c, gin.H{
		"status": d.Status,
	})
}
