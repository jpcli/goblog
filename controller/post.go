package controller

import (
	"goblog/model/request"
	"goblog/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

func PostAdd(c *gin.Context) {
	post := request.Post{}
	err := c.BindJSON(&post)
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

func PostModify(c *gin.Context) {
	post := request.Post{}
	err := c.BindJSON(&post)
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

func PostStatusModify(c *gin.Context) {
	p := request.PostStatusModify{}
	err := c.BindJSON(&p)
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
