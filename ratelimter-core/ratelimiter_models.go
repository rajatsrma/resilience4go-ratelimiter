package ratelimtercore

type Request struct {
	RatelimitKey string
}

type Response struct {
	AvailableRequestQuota int64
	RequestThrottled      bool
}

type RateLimiterMiddlewareConfig struct {
	RateLimitEnabled       bool
	HttpHeaderRateLimitkey string
	Ratelimiter            Ratelimiter
}
