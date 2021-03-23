package handler

import (
	"fmt"
	"goblog/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func MgtViewHandler(c *gin.Context) {
	path := c.Param("path")
	if path == "/" {
		c.HTML(http.StatusOK, "mgt-index.html", nil)
	} else if name := strings.TrimPrefix(path, "/"); utils.Exist(fmt.Sprintf("./view/mgt/includes/%s", name)) {
		c.HTML(http.StatusOK, fmt.Sprintf("mgt-%s", name), nil)
	} else {
		c.String(http.StatusNotFound, "404")
	}
}
