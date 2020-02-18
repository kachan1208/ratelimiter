package ratelimiter

import (
	"errors"
	"time"
)

//Limiter interface describes basic structure of each implemented limiter and allows you to
//create your own realisation without changing code that already use it
type Limiter interface {
	Limit(string, time.Duration, uint32) error
}

var (
	ErrLimitReached = errors.New("rate limit reached")
)
