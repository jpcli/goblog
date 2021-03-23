package config

import (
	"encoding/json"
	"goblog/utils/errors"
	"io/ioutil"
)

type configSet struct {
	db    `json:"db"`
	app   `json:"app"`
	oauth `json:"oauth"`
	jwt   `json:"jwt"`
}

type db struct {
	IP       string `json:"ip"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

type app struct {
	Address string `json:"address"`
	MgtURI  string `json:"mgt_uri"`
}
type jwt struct {
	Issuer string `json:"issuer"`
	Key    string `json:"key"`
}

type oauth struct {
	Admin  string `json:"admin"`
	github `json:"github"`
}

type github struct {
	UseCFWorker  bool   `json:"use-cf-worker"`
	CFWorkerURL  string `json:"cf-worker-url"`
	ClientName   string `json:"client_name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

var config configSet

// LoadConfig which should run before web server load all config defined in 'config.json'.
func LoadConfig() {
	src, err := ioutil.ReadFile("./config.json")
	if err != nil {
		panic(errors.WrapfErrorWithStack(err, "failed to read config"))
	}

	err = json.Unmarshal(src, &config)
	if err != nil {
		panic(errors.WrapfErrorWithStack(err, "failed to load config"))
	}
}

// DBConfig returns database config.
func DBConfig() (ip, port, usr, pwd, db string) {
	return config.db.IP, config.db.Port, config.db.Username, config.db.Password, config.db.Database
}

func AppAdress() string {
	return config.app.Address
}
func AppMgtURI() string {
	return config.app.MgtURI
}

func OauthAdmin() string {
	return config.oauth.Admin
}

func OauthGithubUseCFWorker() bool {
	return config.oauth.UseCFWorker
}
func OauthGithubCFWorkerURL() string {
	return config.oauth.CFWorkerURL
}
func OauthGithubClientName() string {
	return config.oauth.ClientName
}
func OauthGithubClientID() string {
	return config.oauth.ClientID
}
func OauthGithubClientSecret() string {
	return config.oauth.ClientSecret
}

func JwtIssuer() string {
	return config.jwt.Issuer
}
func JwtKey() string {
	return config.jwt.Key
}
