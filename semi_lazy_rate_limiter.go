package ratelimiter

import (
	"fmt"
	"sync"
	"time"
)

type SemiLazyLimiter struct {
	delay    time.Time
	mutex    sync.RWMutex
	locks    map[string]*lock
	gcTicker *time.Ticker
}

const (
	//gcPeriod is a period of time to run garbage collector
	//time is randomly taken
	gcPeriod time.Duration = time.Second * 10
)

func NewSemiLazyLimiter() *SemiLazyLimiter {
	l := SemiLazyLimiter{
		locks:    make(map[string]*lock),
		gcTicker: time.NewTicker(gcPeriod),
	}

	go l.runGCDeamon()

	return &l
}

func (r *SemiLazyLimiter) Limit(key string, period time.Duration, limit uint32) error {
	return r.lock(key, &lock{
		ttl:   time.Now().Add(period).UnixNano(),
		limit: limit,
		count: 1,
	})
}

func (r *SemiLazyLimiter) lock(key string, l *lock) error {
	r.mutex.RLock()
	lock, exists := r.locks[key]
	r.mutex.RUnlock()

	r.mutex.Lock()
	defer r.mutex.Unlock()
	if exists {
		if lock.isTTLReached() {
			r.locks[key] = l
			return nil
		}

		lock.count++
		if lock.count > lock.limit {
			return ErrLimitReached
		}

	} else {
		r.locks[key] = l
	}

	return nil
}

//runGCDeamon used for active cache cleaning, it's trying to copy redis mechanism of cache deactivation
//a bit overenginerred but simply helps to reduce memory usage
func (r *SemiLazyLimiter) runGCDeamon() {
	for _ = range r.gcTicker.C {
		r.clean()
	}
}

func (r *SemiLazyLimiter) clean() {
	//lock Rmutex to read data from storage
	r.mutex.RLock()
	fmt.Println("Start", len(r.locks))
	for key, lock := range r.locks {
		//unlock Rmutex because reading ended
		r.mutex.RUnlock()
		//here we have empty instruction, to allow another threads to catch r/w mutex

		r.mutex.Lock()
		if lock.isTTLReached() {
			delete(r.locks, key)
		}
		r.mutex.Unlock()

		r.mutex.RLock()
	}

	fmt.Println(len(r.locks))
	r.mutex.RUnlock()
}
