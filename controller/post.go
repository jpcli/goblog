package controller

import (
	"fmt"
	"goblog/model/request"
	"goblog/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 新建文章API控制器
func PostAdd(c *gin.Context) {
	post := request.Post{}
	err := c.ShouldBindJSON(&post)
	if err != nil {
		apiErrorInput(c)
		return
	}
	post.ID = 0

	id, err := service.EditPost(&post)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"id": id,
	}, "新建文章成功")
}

// 修改文章API控制器
func PostModify(c *gin.Context) {
	post := request.Post{}
	err := c.ShouldBindJSON(&post)
	if err != nil {
		apiErrorInput(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErrorInput(c)
		return
	}
	post.ID = uint32(id)

	_, err = service.EditPost(&post)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"id": id,
	}, "修改文章成功")
}

// 修改文章状态API控制器
func PostStatusModify(c *gin.Context) {
	p := request.PostStatusModify{}
	err := c.ShouldBindJSON(&p)
	if err != nil {
		apiErrorInput(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErrorInput(c)
		return
	}
	p.ID = uint32(id)

	err = service.ModifyPostStatus(&p)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"status": p.Status,
	}, "修改文章类型成功")
}

// 单个文章获取API控制器
func PostGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErrorInput(c)
		return
	}

	post, err := service.GetPost(uint32(id))
	if err != nil {
		apiFailed(c, err.Error())
		return
	}
	cateID, tagsID, err := service.GetPostCateTags(uint32(id))
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"post":        post,
		"category_id": cateID,
		"tags_id":     tagsID,
	}, "获取文章成功")
}

// 正常文章列表API控制器，参数为?page=&limit=
func PostListNormal(c *gin.Context) {
	p := request.Pager{}
	err := c.ShouldBindQuery(&p)
	if err != nil {
		apiErrorInput(c)
		return
	}

	posts, err := service.ListNormalPost(p.Pi, p.Ps)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	res := make([]gin.H, 0, len(posts))
	for i := range posts {
		cateID, tagsID, err := service.GetPostCateTags(posts[i].Pid)
		if err != nil {
			apiFailed(c, fmt.Sprintf("[ID%d]%s", posts[i].Pid, err.Error()))
			return
		}
		res = append(res, gin.H{
			"post":        &posts[i],
			"category_id": cateID,
			"tags_id":     tagsID,
		})
	}

	apiOK(c, gin.H{
		"list": res,
	})
}
