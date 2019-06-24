package RedisClient

import (
	"github.com/go-redis/redis"
	"os"
)

func NewRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password:  os.Getenv("REDIS_PASSWORD"),
		DB: 0,
	})
}
