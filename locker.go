package ratelimiter

import "time"

type lock struct {
	ttl   int64
	limit uint32
	count uint32
}

func (l *lock) isTTLReached() bool {
	return l.ttl < time.Now().UnixNano()
}
