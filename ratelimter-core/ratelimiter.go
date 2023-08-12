package ratelimtercore

type Ratelimiter interface {
	checkForRateLimit(Request) Response
}
