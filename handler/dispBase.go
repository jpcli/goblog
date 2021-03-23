package handler

import (
	"fmt"
	"goblog/service"
	"goblog/utils/errors"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func HTMLNotFound(c *gin.Context, err error) {
	c.HTML(http.StatusNotFound, "disp-404.html", gin.H{})
	//TODO:写日志
	log.Println(errors.SprintError(err))
}

func HTMLServerError(c *gin.Context, err error) {
	c.HTML(http.StatusInternalServerError, "disp-500.html", gin.H{})
	//TODO:写日志
	log.Println(errors.SprintError(err))
}

func handleServiceError(c *gin.Context, e *service.DispError) {
	switch e.GetErrorType() {
	case service.ErrorRequest:
		HTMLNotFound(c, e.GetError())
	case service.ErrorServer:
		HTMLServerError(c, e.GetError())
	}
}

type tkd struct {
	Title       string
	Description string
	Keywords    string
}

func HTMLOK(c *gin.Context, name string, data gin.H, tkd *tkd) {
	data["title"] = tkd.Title
	data["description"] = tkd.Description
	data["keywords"] = tkd.Keywords

	s := service.NewService()
	data["categories"] = s.GetAllDispCategories()
	data["tags"] = s.GetAllDispTags()
	count, _ := s.CountPost()
	data["postNum"] = count
	data["lastModified"] = s.GetLastModified()

	c.HTML(http.StatusOK, name, data)
}

func getNumFromRegexp(exp string, s string) (uint32, error) {
	comp, err := regexp.Compile(exp)
	if err != nil {
		return 0, fmt.Errorf("failed to compile regexp")
	}
	d := comp.FindStringSubmatch(s)
	if d == nil {
		return 0, fmt.Errorf("no number matched with the given expression")
	}
	num, err := cast.ToUint32E(d[1])
	if err != nil {
		return 0, err
	}
	return num, nil
}
