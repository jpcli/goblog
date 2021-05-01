package controller

import (
	"goblog/service"
	"strings"

	"github.com/gin-gonic/gin"
)

// 上传文件API控制器
func Upload(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		apiFailed(c, "接受文件失败")
		return
	}

	savePath, err := service.AddAttachment(file.Filename)
	if err != nil {
		apiFailed(c, "添加文件失败")
	}

	err = c.SaveUploadedFile(file, savePath)
	if err != nil {
		apiFailed(c, "保存文件失败")
	}

	apiOK(c, gin.H{
		"uri": strings.Replace(savePath, "./", "/", 1),
	}, "保存文件成功")
}
