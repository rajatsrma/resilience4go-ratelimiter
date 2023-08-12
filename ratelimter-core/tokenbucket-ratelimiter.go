package ratelimtercore

import (
	"context"
	"log"
	"math"
	"strconv"
	"time"
)

type TokenBucketRateLimiter struct {
	BucketCapacity  int
	RefillFrequency int
}

func (t TokenBucketRateLimiter) CheckForRateLimit(request Request) Response {
	ctx := context.Background()

	if redisClient == nil {
		log.Default().Println("[GetRatelimiterRedis]: issue in creating connection with redis")
		return Response{RequestThrottled: false}
	}

	result, err := redisClient.HMGet(ctx, request.RatelimitKey, "tokens", "updateTime").Result()
	if err != nil {
		return Response{RequestThrottled: true}
	}

	if result[0] == nil || result[1] == nil {
		tokens := t.BucketCapacity - 1
		updateTime := time.Now().UnixMilli()

		_, err := redisClient.HSet(ctx, request.RatelimitKey, "tokens", tokens, "updateTime", updateTime).Result()
		if err != nil {
			log.Default().Println("[KEY_NOT_EXISTS]: error while setting the token value", err)
		}

		return Response{RequestThrottled: false, AvailableRequestQuota: int64(tokens)}
	} else {
		lastUpdateTimeMs, _ := strconv.Atoi(result[1].(string))
		currentTime := time.Now().UnixMilli()

		pendingTokens, _ := strconv.Atoi(result[0].(string))
		newEarnedTokens := ((currentTime - int64(lastUpdateTimeMs)) / 50000) * int64(t.RefillFrequency)
		totalTokens := math.Min(float64(int64(pendingTokens)+newEarnedTokens), float64(t.BucketCapacity)) - 1

		if totalTokens < 0 {
			return Response{RequestThrottled: true, AvailableRequestQuota: int64(0)}
		}

		_, err := redisClient.HSet(ctx, request.RatelimitKey, "tokens", totalTokens, "updateTime", currentTime).Result()
		if err != nil {
			log.Default().Println("[KEY_EXISTS]: error while setting the token value", err)
		}
		return Response{RequestThrottled: false, AvailableRequestQuota: int64(totalTokens)}
	}
}
