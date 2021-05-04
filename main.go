package main

import (
	"goblog/dao"
	"goblog/router"
	"goblog/utils/config"
)

func main() {
	config.LoadConfig()
	dao.OpenDatabase(config.MysqlIP(), config.MysqlPort(), config.MysqlUser(), config.MysqlPwd(), config.MysqlDB())
	router.AppRun(config.AppAddr())
}
