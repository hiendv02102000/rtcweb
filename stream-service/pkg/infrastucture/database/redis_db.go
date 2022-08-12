package db

import (
	goRedis "github.com/go-redis/redis/v8"
)

var RedisPool = ConnectRedisPool(0)

func ConnectRedisPool(db int) *goRedis.Client {
	rdb := goRedis.NewClient(&goRedis.Options{
		Addr:     "redis-app:6379",
		Password: "password",
		DB:       db,
		PoolSize: 1000,
	})
	return rdb
}
