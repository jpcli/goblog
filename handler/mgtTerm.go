package handler

import (
	"fmt"
	"goblog/service"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

func GetTermAPI(c *gin.Context) {
	tid, err := cast.ToUint32E(c.DefaultQuery("id", "0"))
	if err != nil {
		APIError(c, fmt.Sprintf("unexpected id value: %s", err))
		return
	}

	s := service.NewService()

	term, err := s.GetMgtTermByTid(tid)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get term: %s", err))
		return
	}

	APIOK(c, term)
}

func GetTermListAPI(c *gin.Context) {
	termType := c.Query("type")

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

	termList, err := s.GetMgtTermList(termType, page, limit)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get term list: %s", err))
		return
	}

	count, err := s.CountTerm(termType)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to count term: %s", err))
		return
	}

	APIOK(c, gin.H{
		"count": count,
		"list":  termList,
	})

}

func EditTermAPI(c *gin.Context) {
	tid, err := cast.ToUint32E(c.DefaultQuery("id", "0"))
	if err != nil {
		APIError(c, fmt.Sprintf("unexpected id value: %s", err))
		return
	}

	var data service.MgtTerm
	err = c.ShouldBindJSON(&data)
	if err != nil {
		APIError(c, fmt.Sprintf("miss term info: %s", err))
		return
	}

	s := service.NewService()

	id, err := s.EditTerm(tid, &data)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to edit term: %s", err))
		return
	}

	APIOK(c, gin.H{
		"id": id,
	})
}
