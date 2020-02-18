package ratelimiter

import (
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestSemiLazyRateLimiterConstructor(t *testing.T) {
	limiter := NewSemiLazyLimiter()
	assert.NotNil(t, limiter)
}

func TestSemiLazyRateLimiterLimitLockedSuccessfuly(t *testing.T) {
	limiter := NewSemiLazyLimiter()
	assert.NotNil(t, limiter)

	testFunc := func() error {
		return limiter.Limit(t.Name(), time.Millisecond*100, 1)
	}

	emptyErr := testFunc()
	nonEmptyErr := testFunc()

	assert.NoError(t, emptyErr)
	assert.Error(t, nonEmptyErr)
}

func TestSemiLazyRateLimiterLimitTTLReachedAndUpdatedSuccess(t *testing.T) {
	limiter := NewSemiLazyLimiter()
	assert.NotNil(t, limiter)

	testFunc := func() error {
		return limiter.Limit(t.Name(), time.Millisecond*50, 1)
	}

	err := testFunc()
	assert.NoError(t, err)

	time.Sleep(time.Millisecond * 60)
	err = testFunc()
	assert.NoError(t, err)
}

func TestSemiLazyRateLimiterLimitDifferentKeys(t *testing.T) {
	limiter := NewSemiLazyLimiter()
	assert.NotNil(t, limiter)

	testFunc := func(key string) error {
		return limiter.Limit(key, time.Millisecond*50, 1)
	}

	assert.NoError(t, testFunc("1"))
	assert.NoError(t, testFunc("2"))
	assert.NoError(t, testFunc("3"))
}

func TestLockIsTTLReachedSuccess(t *testing.T) {
	l := lock{
		ttl: time.Now().Add(time.Second).UnixNano(),
	}

	assert.False(t, l.isTTLReached())
}

func TestLockIsTTLReachedFail(t *testing.T) {
	l := lock{
		ttl: time.Now().Add(-time.Second).UnixNano(),
	}

	assert.True(t, l.isTTLReached())
}

func TestSemiLazyRateLimiterLimitConcurrent(t *testing.T) {
	limiter := NewSemiLazyLimiter()
	testFunc := func() error {
		return limiter.Limit("1", time.Second*5, 100)
	}

	var wg errgroup.Group
	for i := 0; i < 100; i++ {
		wg.Go(testFunc)
	}

	err := wg.Wait()

	assert.NoError(t, err)
	assert.Error(t, testFunc())
}

func BenchmarkSemiLazyRateLimiterLimitDifferentKeys(b *testing.B) {
	limiter := NewSemiLazyLimiter()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		limiter.Limit(strconv.Itoa(b.N), time.Second*1, 1)
	}
}
