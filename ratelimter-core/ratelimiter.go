package ratelimtercore

type Ratelimiter interface {
	CheckForRateLimit(Request) Response
}
