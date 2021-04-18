package main

import (
	"goblog/cache"
	"goblog/repository"
	"goblog/router"
	"goblog/utils/config"
	"goblog/utils/option"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	config.LoadConfig()
	repository.OpenDatabase(config.DBConfig())
	cache.InitCache(config.RedisConfig())
	option.InitOption()
	router.AppStart(config.AppAddress())
}
