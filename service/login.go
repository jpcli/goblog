package service

import (
	"encoding/json"
	"fmt"
	"goblog/utils/config"
	"goblog/utils/jwt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// 直接请求github服务器
func GithubLoginDirectly(code string) (map[string]interface{}, error) {
	// 通过code获取access_token
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s",
		config.OauthGithubClientID(), config.OauthGithubClientSecret(), code,
	)
	data, err := githubGetJSON(url, "POST", "", "")
	if err != nil {
		return nil, err
	} else if errorString, ok := data["error"]; ok {
		return nil, detectGithubError(errorString.(string))
	}

	// 通过access_token获取用户信息
	data, err = githubGetJSON("https://api.github.com/user", "GET", config.OauthGithubClientName(), data["access_token"].(string))
	return data, err
}

// 通过cfworker代理请求github服务器
func GithubLoginWithCFWorker(code string) (map[string]interface{}, error) {
	data, err := githubGetJSON(fmt.Sprintf("%s?code=%s", config.OauthGithubCFWorkerURL(), code), "GET", "", "")
	if err != nil {
		return nil, err
	} else if errorString, ok := data["error"]; ok {
		return nil, detectGithubError(errorString.(string))
	}
	return data, err
}

// 生成github的JWT
func GenerateGithubJwt(data map[string]interface{}) (string, error) {
	claims := &jwt.CustomClaims{}
	if user, ok := data["login"]; ok && user.(string) != "" {
		claims.User = user.(string)
	} else {
		return "", fmt.Errorf("github用户名丢失")
	}

	token, err := jwt.NewJWT(claims)
	if err != nil {
		return "", err
	}
	return token, nil
}

// 侦测github返回的错误类型
func detectGithubError(errMsg string) error {
	switch errMsg {
	case "bad_verification_code":
		return fmt.Errorf("code不正确或已过期，请重试")
	case "incorrect_client_credentials":
		// TODO 写日志
		return fmt.Errorf("github应用配置失败，请联系管理员")
	default:
		// TODO 写日志
		return fmt.Errorf("发生未知错误，请联系管理员")
	}
}

// 请求github服务器，ua、accessToken可为空，返回反序列化后的json数据或者错误
func githubGetJSON(url, method, ua, accessToken string) (map[string]interface{}, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	// 设置请求头
	request, err := http.NewRequest(method, url, strings.NewReader(""))
	if err != nil {
		// TODO 写日志
		return nil, fmt.Errorf("服务器发生错误")
	}
	request.Header.Set("accept", "application/json")
	if ua != "" {
		request.Header.Set("User-Agent", ua)
	}
	if accessToken != "" {
		request.Header.Set("Authorization", fmt.Sprintf("bearer %s", accessToken))
	}

	// 发送请求
	response, err := client.Do(request)
	if err != nil {
		if strings.Contains(err.Error(), "Client.Timeout exceeded") {
			return nil, fmt.Errorf("请求github服务器超时")
		} else {
			// TODO 写日志
			return nil, fmt.Errorf("请求github服务器失败")
		}
	}

	// 解析数据
	res, _ := ioutil.ReadAll(response.Body)
	var data map[string]interface{}
	err = json.Unmarshal(res, &data)
	if err != nil {
		// TODO 写日志
		return nil, fmt.Errorf("解析github返回数据失败")
	}

	// 关闭连接、返回
	response.Body.Close()
	return data, nil
}
