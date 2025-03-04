package ioc

import (
	"github.com/redis/go-redis/v9"
	"webookProgram/webook/config"
)

func InitCache() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.Config.Redis.Addr,
		Password: "",
	})
}
