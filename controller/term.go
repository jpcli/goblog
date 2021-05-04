package controller

import (
	"goblog/model"
	"goblog/model/request"
	"goblog/service"
	"math"
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

	count, err := service.CountTermByType(t.TermType)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	apiOK(c, gin.H{
		"list":  terms,
		"count": count,
	}, "获取项列表成功")
}

// 获取所有分类、标签API控制器
func GetAllCategoryTags(c *gin.Context) {
	cates, err := service.ListTerm(uint8(model.TERM_TYPE_CATEGORY), 1, math.MaxUint32)
	if err != nil {
		apiFailed(c, "获取所有分类失败")
		return
	}
	tags, err := service.ListTerm(uint8(model.TERM_TYPE_TAG), 1, math.MaxUint32)
	if err != nil {
		apiFailed(c, "获取所有标签失败")
		return
	}

	apiOK(c, gin.H{
		"categories": cates,
		"tags":       tags,
	}, "获取所有分类、标签成功")
}
