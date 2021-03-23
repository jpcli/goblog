package handler

import (
	"fmt"
	"goblog/service"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadAPI(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get file: %s", err))
	}

	s := service.NewService()
	imgPath, err := s.Upload(file.Filename)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to upload file: %s", err))
	}

	err = c.SaveUploadedFile(file, imgPath)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to sava file: %s", err))
	}

	APIOK(c, gin.H{
		"uri": strings.Replace(imgPath, "./", "/", 1),
	})
}
