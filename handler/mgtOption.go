package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goblog/service"
)

func GetWebsiteBasicOptionAPI(c *gin.Context) {
	s := service.NewService()
	options, err := s.GetWebsiteBasicOption()
	if err != nil {
		APIError(c, fmt.Sprintf("failed to get options: %s", err))
		return
	}

	APIOK(c, options)
}

func ModifyWebsiteBasicOptionAPI(c *gin.Context) {
	type websiteBasicOption struct {
		WebsiteURL    string `json:"websiteURL"`
		WebsiteName   string `json:"websiteName"`
		EachPageLimit uint32 `json:"eachPageLimit"`
		PageNavLimit  uint32 `json:"pageNavLimit"`
	}

	d := websiteBasicOption{
		EachPageLimit: 10,
		PageNavLimit:  7,
	}
	err := c.BindJSON(&d)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to bind request data: %s", err))
		return
	}

	s := service.NewService()
	err = s.SetWebsiteBasicOption(d.WebsiteURL, d.WebsiteName, d.EachPageLimit, d.PageNavLimit)
	if err != nil {
		APIError(c, fmt.Sprintf("failed to modify options: %s", err))
		return
	}

	APIOK(c, gin.H{})
}
