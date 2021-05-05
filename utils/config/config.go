package config

import (
	"encoding/json"
	"io/ioutil"
)

type configSet struct {
	mysql `json:"mysql"`
	redis `json:"redis"`
	oauth `json:"oauth"`
	jwt   `json:"jwt"`
	app   `json:"app"`
}

var config configSet

// 加载配置，在程序开始执行该命令
func LoadConfig() {
	src, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic("读取应用配置config.json失败")
	}

	err = json.Unmarshal(src, &config)
	if err != nil {
		panic("解析应用配置config.json失败，请检查语法是否正确")
	}
}

// mysql配置
type mysql struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pwd  string `json:"pwd"`
	DB   string `json:"db"`
}

func MysqlIP() string {
	return config.mysql.IP
}
func MysqlPort() int {
	return config.mysql.Port
}
func MysqlUser() string {
	return config.mysql.User
}
func MysqlPwd() string {
	return config.mysql.Pwd
}
func MysqlDB() string {
	return config.mysql.DB
}

// redis配置
type redis struct {
	IP   string `json:"ip"`
	Port int    `json:"port"`
	Pwd  string `json:"pwd"`
	DB   int    `json:"db"`
}

func RedisIP() string {
	return config.redis.IP
}
func RedisPort() int {
	return config.redis.Port
}
func RedisPwd() string {
	return config.redis.Pwd
}
func RedisDB() int {
	return config.redis.DB
}

// oauth配置
type oauth struct {
	github `json:"github"`
}
type github struct {
	AdminUser    string `json:"admin_user"`
	UseCFWorker  bool   `json:"use_cf_worker"`
	CFWorkerURL  string `json:"cf_worker_url"`
	ClientName   string `json:"client_name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func OauthGithubAdminUser() string {
	return config.oauth.github.AdminUser
}
func OauthGithubUseCFWorker() bool {
	return config.oauth.github.UseCFWorker
}
func OauthGithubCFWorkerURL() string {
	return config.oauth.github.CFWorkerURL
}
func OauthGithubClientName() string {
	return config.oauth.github.ClientName
}
func OauthGithubClientID() string {
	return config.oauth.github.ClientID
}
func OauthGithubClientSecret() string {
	return config.oauth.github.ClientSecret
}

// jwt配置
type jwt struct {
	Issuer string `json:"issuer"`
	Key    string `json:"key"`
}

func JwtIssuer() string {
	return config.jwt.Issuer
}
func JwtKey() string {
	return config.jwt.Key
}

// app应用运行配置
type app struct {
	Addr              string `json:"addr"`
	AdminSafetyFactor string `json:"admin_safety_factor"`
	LogFile           string `json:"log_file"`
}

func AppAddr() string {
	return config.app.Addr
}
func AdminSafetyFactor() string {
	return config.app.AdminSafetyFactor
}
func AppLogFile() string {
	return config.app.LogFile
}
