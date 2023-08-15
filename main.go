package main

import (
	"github.com/gin-gonic/gin"
	ratelimtercore "github.com/rajatsrma/resilience4go-ratelimiter/ratelimter-core"
)

const (
	redisConnectionString = "redis://localhost:6379/0"
)

func main() {
	r := gin.Default()
	r.Use(ratelimtercore.SetRateLimiterRedisMiddleware(redisConnectionString))

	// Sliding window counter rate limiting
	r.GET(
		"/",
		ratelimtercore.RateLimiterMiddleware(
			ratelimtercore.RateLimiterMiddlewareConfig{
				RateLimitEnabled:       true,
				HttpHeaderRateLimitkey: "rate-limit-key",
				Ratelimiter: ratelimtercore.SlidingWindowCounterRatelimiter{
					TimeWindowInSeconds:  120,
					AllowedRequestsCount: 10,
				},
			},
		),
		handleRequest,
	)
	// token bucket rate limiting
	r.GET(
		"/token",
		ratelimtercore.RateLimiterMiddleware(
			ratelimtercore.RateLimiterMiddlewareConfig{
				RateLimitEnabled:       true,
				HttpHeaderRateLimitkey: "rate-limit-key",
				Ratelimiter: ratelimtercore.TokenBucketRatelimiter{
					BucketCapacity:  6,
					RefillFrequency: 10,
				},
			},
		),
		handleRequest,
	)
	//  No rate limiting
	r.GET(
		"/other",
		ratelimtercore.RateLimiterMiddleware(
			ratelimtercore.RateLimiterMiddlewareConfig{RateLimitEnabled: false},
		),
		handleRequest,
	)

	r.Run(":8080")
}

func handleRequest(c *gin.Context) {
	c.String(200, "Hello, Gin Server!")
}
