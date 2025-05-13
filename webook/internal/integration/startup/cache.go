package startup

import (
	"github.com/redis/go-redis/v9"
)

var redisClient redis.Cmdable

func InitCache() redis.Cmdable {
	if redisClient == nil {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
		})
	}
	return redisClient
}
