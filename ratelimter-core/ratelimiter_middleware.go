package ratelimtercore

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func SetRateLimiterRedisMiddleware(redisConnectionString string) gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("Middleware: set redis")
		SetRatelimiterRedis(redisConnectionString)
		fmt.Println("Middleware: set redis done")
		c.Next()
	}
}

func RateLimiterMiddleware(config RateLimiterMiddlewareConfig) gin.HandlerFunc {
	tokenBucketRateLimiter := TokenBucketRatelimiter{BucketCapacity: config.MaxRps, RefillFrequency: config.RefillFrequency}
	return func(c *gin.Context) {
		if config.RateLimitEnabled {
			fmt.Println("Middleware: Before request")

			rateLimitKey := string(c.GetHeader(config.HttpHeaderRateLimitkey))
			if rateLimitKey == "" {
				c.AbortWithStatusJSON(400, gin.H{"error": "http headers not properly configured for request"})
			}

			rateLimitResponse := tokenBucketRateLimiter.CheckForRateLimit(
				Request{RatelimitKey: rateLimitKey},
			)
			fmt.Println("ratelimitResponse", rateLimitResponse)
			if rateLimitResponse.RequestThrottled {
				c.AbortWithStatusJSON(429, gin.H{"ratelimitResponse": rateLimitResponse})
			} else {
				c.Next()
			}
		} else {
			c.Next()
		}
	}
}