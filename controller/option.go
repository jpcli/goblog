package controller

import (
	"goblog/model/request"
	"goblog/service"

	"github.com/gin-gonic/gin"
)

// 获取基本选项API控制器
func BaseOptionGet(c *gin.Context) {
	baseOption, err := service.GetBaseOption()
	if err != nil {
		apiFailed(c, "获取基本选项失败")
		return
	}

	res := gin.H{}
	for name, val := range baseOption {
		res[name] = val
	}
	apiOK(c, res, "获取基本选项成功")
}

// 更新基本选项API控制器
func BaseOptionSet(c *gin.Context) {
	opt := request.BaseOption{
		PageSize:    10,
		PageNavSize: 7,
	}
	err := c.ShouldBindJSON(&opt)
	if err != nil || opt.PageSize <= 0 || opt.PageNavSize <= 0 {
		apiErrorInput(c)
		return
	}

	err = service.SetBaseOption(&opt)
	if err != nil {
		apiFailed(c, "设置基本选项失败")
		return
	}

	apiOK(c, gin.H{}, "设置基本选项成功")
}
