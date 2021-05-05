package main

import (
	"goblog/dao"
	"goblog/router"
	"goblog/utils/config"
)

func main() {
	config.LoadConfig()
	dao.OpenDatabase()
	router.AppRun()
}
