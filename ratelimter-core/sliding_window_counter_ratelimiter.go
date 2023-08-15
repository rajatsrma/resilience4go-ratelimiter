package ratelimtercore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	redis "github.com/redis/go-redis/v9"
)

type SlidingWindowCounterRatelimiter struct {
	TimeWindowInSeconds  int64
	AllowedRequestsCount int64
}

func (s SlidingWindowCounterRatelimiter) CheckForRateLimit(request Request) Response {
	ctx := context.Background()

	if redisClient == nil {
		log.Default().Println("[GetRatelimiterRedis]: issue in creating connection with redis")
		return Response{RequestThrottled: false}
	}

	timestampSlice, err := getAllowedRequestData(ctx, request.RatelimitKey)
	if err != nil {
		return Response{RequestThrottled: true}
	}

	currentTime, filteredTimestamps := getTimestampsInTimeWindow(timestampSlice, s.TimeWindowInSeconds)
	filteredTimestamps = append(filteredTimestamps, strconv.FormatInt(currentTime, 10))

	seterr := setTimestampsInRedis(ctx, request.RatelimitKey, filteredTimestamps)
	if seterr != nil {
		return Response{RequestThrottled: true}
	}

	if len(filteredTimestamps) > int(s.AllowedRequestsCount) {
		return Response{RequestThrottled: true}
	} else {
		return Response{RequestThrottled: false, AvailableRequestQuota: int64(int(s.AllowedRequestsCount) - len(filteredTimestamps))}
	}
}

func getAllowedRequestData(ctx context.Context, redisKey string) ([]string, error) {
	result, err := redisClient.Get(ctx, redisKey).Result()
	if err == redis.Nil {
		return []string{}, nil
	} else if err != nil {
		fmt.Println("Error get value for key:", err)
		return nil, err
	}

	var timestampSlice []string
	unmarshalErr := json.Unmarshal([]byte(result), &timestampSlice)
	if unmarshalErr != nil {
		fmt.Println("Error:", err)
		return nil, unmarshalErr
	}
	return timestampSlice, nil
}

func getTimestampsInTimeWindow(timestamps []string, timeWindowInSeconds int64) (int64, []string) {
	currentTime := time.Now().UnixMilli()
	timeWindowStart := currentTime - (timeWindowInSeconds * 1000)

	// check number of request within these limits
	filteredTimestamps := filter(timestamps, func(timestamp string) bool {
		timestampInt, err := strconv.ParseInt(timestamp, 10, 64)
		if err != nil {
			log.Default().Println("issue in converting timestamp string to int64")
			return true
		}
		return int64(timestampInt) >= timeWindowStart && int64(timestampInt) <= currentTime
	})
	return currentTime, filteredTimestamps
}

func setTimestampsInRedis(ctx context.Context, ratelimitKey string, timestamps []string) error {
	serialisedTimestamps, serialseErr := json.Marshal(timestamps)
	if serialseErr != nil {
		fmt.Println("Error serializing list to JSON:", serialseErr)
		return serialseErr
	}
	_, err := redisClient.Set(ctx, ratelimitKey, serialisedTimestamps, 0).Result()
	if err != nil {
		return err
	}
	return nil
}
