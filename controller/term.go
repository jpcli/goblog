package controller

import (
	"goblog/model/request"
	"goblog/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 新建项API控制器
func TermAdd(c *gin.Context) {
	term := request.Term{}
	err := c.ShouldBindJSON(&term)
	if err != nil {
		apiErrorInput(c)
		return
	}
	term.ID = 0

	id, err := service.EditTerm(&term)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"id": id,
	}, "新建项成功")
}

// 修改项API控制器
func TermModify(c *gin.Context) {
	term := request.Term{}
	err := c.ShouldBindJSON(&term)
	if err != nil {
		apiErrorInput(c)
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErrorInput(c)
		return
	}
	term.ID = uint32(id)

	_, err = service.EditTerm(&term)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"id": term.ID,
	}, "修改项成功")
}

// 获取单个项API控制器
func TermGet(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		apiErrorInput(c)
		return
	}

	term, err := service.GetTerm(uint32(id))
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"term": term,
	}, "获取项成功")

}

// 获取项列表API控制器，参数为?page=&limit=
func TermList(c *gin.Context) {
	t := request.TermList{}
	err := c.ShouldBindQuery(&t)
	if err != nil || t.Pi <= 0 || t.Ps <= 0 {
		apiErrorInput(c)
		return
	}

	terms, err := service.ListTerm(t.TermType, t.Pi, t.Ps)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"list": terms,
	})
}
