package cache

import (
	"fmt"
	"github.com/go-redis/redis"
)

var redisDB *redis.Client

func InitCache(ip, port, pwd string, db int) {
	redisDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", ip, port),
		Password: pwd,
		DB:       db,
	})
	_, err := redisDB.Ping().Result()
	if err != nil {
		panic("failed to connect to redis database")
	}
}

type Cache struct {
	redis *redis.Client
}

var NoResultErr = fmt.Errorf("cache no result")

func NewCache() *Cache {
	return &Cache{redis: redisDB}
}

func (c *Cache) GetOption(key string) (string, error) {
	value, err := c.redis.HGet("options", key).Result()
	if err == redis.Nil {
		return "", NoResultErr
	} else if err != nil {
		return "", err
	} else {
		return value, nil
	}
}

func (c *Cache) SetOption(key, value string) error {
	err := c.redis.HSet("options", key, value).Err()
	return err
}

func (c *Cache) ClearOption() error {
	return c.redis.Del("options").Err()
}
