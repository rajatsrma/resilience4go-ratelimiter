package ratelimtercore

import (
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient     *redis.Client
	redisClientLock sync.Mutex
)

func SetRatelimiterRedis(connectionString string) {
	// create connection only if doesn't exists
	if redisClient != nil {
		return
	}
	redisClientLock.Lock()
	defer redisClientLock.Unlock()

	if redisClient != nil {
		return
	}

	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		fmt.Println("error setting connection with rate limiter redis")
	}

	redisClient = redis.NewClient(opt)
}
