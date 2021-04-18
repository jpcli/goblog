package option

import (
	"github.com/spf13/cast"
	"goblog/cache"
	"goblog/repository"
)

var r *repository.Repository
var c *cache.Cache

func InitOption() {
	r = repository.NewRepository()
	c = cache.NewCache()
}

func getOption(key string) string {
	val, err := c.GetOption(key)
	if err != nil {
		val = r.GetOption(key)
		c.SetOption(key, val)
	}
	return val
}

func GetEachPageLimit() uint32 {
	return cast.ToUint32(getOption("eachPageLimit"))
}

func GetPageNavLimit() uint32 {
	return cast.ToUint32(getOption("pageNavLimit"))
}

func GetWebsiteURL() string {
	return getOption("websiteURL")
}

func GetWebsiteName() string {
	return getOption("websiteName")
}
