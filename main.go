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

	r.GET(
		"/",
		ratelimtercore.RateLimiterMiddleware(
			ratelimtercore.RateLimiterMiddlewareConfig{RateLimitEnabled: true, MaxRps: 5, RefillFrequency: 5, HttpHeaderRateLimitkey: "rate-limit-key"},
		),
		handleRequest,
	)
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
