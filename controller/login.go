package controller

import (
	"fmt"
	"goblog/service"
	"goblog/utils/config"
	"path"

	"github.com/gin-gonic/gin"
)

// Github登录API控制器
func GithubLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		apiErrorInput(c)
		return
	}

	var githubData map[string]interface{}
	var err error
	if config.OauthGithubUseCFWorker() {
		githubData, err = service.GithubLoginWithCFWorker(code)
	} else {
		githubData, err = service.GithubLoginDirectly(code)
	}
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	token, err := service.GenerateGithubJwt(githubData)
	if err != nil {
		apiFailed(c, err.Error())
		return
	}

	if githubData["login"].(string) == config.OauthGithubAdminUser() {
		apiOK(c, gin.H{
			"jwt": token,
			"to":  fmt.Sprintf("%s/view/", path.Join("/admin", config.AdminSafetyFactor())),
		}, "登录成功")
	} else {
		apiOK(c, gin.H{
			"jwt": token,
		}, "登录成功")
	}

}
