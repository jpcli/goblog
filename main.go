package main

import (
	"goblog/repository"
	"goblog/router"
	"goblog/utils/config"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadConfig()
	repository.OpenDatabase(config.DBConfig())
	router.AppStart(config.AppAdress())
}
