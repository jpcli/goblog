package handler

import (
	"encoding/json"
	"fmt"
	"goblog/utils/config"
	"goblog/utils/jwt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func oauthSuccess(c *gin.Context, jwt, user string) {
	to := ""
	if user == config.OauthAdmin() {
		to = config.AppMgtURI() + "/view/"
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "登录成功",
		"jwt":  jwt,
		"to":   to,
	})
}
func oauthFail(c *gin.Context, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code": -1,
		"msg":  fmt.Sprintf("登录失败：%s", msg),
	})
}

func GithubOauth(c *gin.Context) {
	if config.OauthGithubUseCFWorker() {
		githubOauthUseCFWorker(c)
	} else {
		githubOauthDirect(c)
	}
}

// githubOauthDirect sends the request to github by this app itsels.
// This always used when this app is running in machine out of China.
func githubOauthDirect(c *gin.Context) {
	code := c.Query("code")

	// get accessToken
	client := &http.Client{Timeout: 15 * time.Second}
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		config.OauthGithubClientID(),
		config.OauthGithubClientSecret(),
		code,
	)
	request, err := http.NewRequest("POST", url, strings.NewReader(""))
	if err != nil {
		oauthFail(c, "服务器错误，请重试")
		return
	}
	request.Header.Set("accept", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			oauthFail(c, "请求授权码超时，请重试")
		} else {
			// TODO:写日志
			oauthFail(c, "请求授权码失败，请重试")
		}
		return
	}
	// get access token
	res, _ := ioutil.ReadAll(resp.Body)
	var data map[string]interface{}
	_ = json.Unmarshal([]byte(res), &data)
	if errorString, ok := data["error"]; ok {
		oauthFail(c, fmt.Sprintf("[%s] %s", errorString, strings.Replace(data["error_description"].(string), "+", " ", -1)))
		return
	}
	accessToken := data["access_token"]
	resp.Body.Close()

	// get user info
	client = &http.Client{Timeout: 15 * time.Second}
	request, _ = http.NewRequest("GET", "https://api.github.com/user", strings.NewReader(""))
	request.Header.Set("accept", "application/json")
	request.Header.Set("User-Agent", config.OauthGithubClientName())
	request.Header.Set("Authorization", fmt.Sprintf("bearer %s", accessToken))
	resp, err = client.Do(request)
	if err != nil {
		oauthFail(c, "请求用户数据失败，请重试")
		return
	}
	res, _ = ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	generateGithubJwt(c, res)
}

// githubOauthUseCFWorker sends the request to cf worker (the code of cf-worker.js in this folder) proxy to get user info.
// By using this func, should deplay the cf-worker first.
// This always used when this app is running in machine in China.
func githubOauthUseCFWorker(c *gin.Context) {
	code := c.Query("code")

	client := &http.Client{Timeout: 15 * time.Second}
	request, _ := http.NewRequest("GET", fmt.Sprintf("%s?code=%s", config.OauthGithubCFWorkerURL(), code), strings.NewReader(""))
	resp, err := client.Do(request)
	if err != nil {
		oauthFail(c, "请求用户数据失败，请重试")
		return
	}
	res, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	generateGithubJwt(c, res)
}

func generateGithubJwt(c *gin.Context, data []byte) {
	d := struct {
		Login string `json:"login"`
		ID    uint32 `json:"id"`
		Name  string `json:"name"`
	}{}
	_ = json.Unmarshal(data, &d)
	token, err := jwt.NewJWT(&jwt.CustomClaims{
		GithubID:   d.ID,
		GithubUser: d.Login,
		GithubName: d.Name,
	})
	if err != nil {
		oauthFail(c, "解析用户数据失败，请重试")
		return
	}

	oauthSuccess(c, token, d.Login)
}
