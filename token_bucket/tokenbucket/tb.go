package tokenbucket

import (
	"sync"
	"time"
)

type TokenBucket struct {
	capacity   int        // 버킷이 담을 수 있는 최대 토큰 수
	rate       int        // 초당 추가할 토큰 수 (초당 토큰수)
	tokens     int        // 버킷의 현재 토큰 수
	lastRefill time.Time  // 마지막 토큰 리필의 타임스탬프
	mutex      sync.Mutex // 동시 액세스를 보호하기 위한 뮤텍스
}

func NewTokenBucket(capacity, rate int) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		rate:       rate,
		tokens:     capacity,
		lastRefill: time.Now(),
	}
}

func (tb *TokenBucket) refill() {
	// 마지막 리필 이후 경과된 시간을 기준으로 버킷에 토큰을 추가

	now := time.Now()
	elapsed := now.Sub(tb.lastRefill)

	// 추가해야 하는 토큰 수 계산
	tokensToAdd := int(elapsed.Seconds() * float64(tb.rate))

	if tokensToAdd > 0 {
		tb.lastRefill = now
		// 토큰을 추가하되 버킷의 용량을 초과하지 않도록.
		tb.tokens = min(tb.tokens+tokensToAdd, tb.capacity)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (tb *TokenBucket) Take(tokens int) bool {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	// 먼저 경과된 시간을 기준으로 버킷을 토큰으로 리필
	tb.refill()

	// 토큰이 있는지 체크
	if tb.tokens >= tokens {
		tb.tokens -= tokens
		return true
	}

	// 사용 가능한 토큰이 충분하지 않음
	return false
}

func (tb *TokenBucket) TakeWithTimeout(tokens int, timeout time.Duration) bool {
	// 기다릴 수 있는 가장 빠른 시간 계산
	deadline := time.Now().Add(timeout)

	for {
		tb.mutex.Lock()

		// bucket refill
		tb.refill()

		// 토큰이 충분한지 체크
		if tb.tokens >= tokens {
			tb.tokens -= tokens
			tb.mutex.Unlock()
			return true
		}

		// 더 많은 토큰을 기다려야 하는 시간 계산
		tokensNeeded := tokens - tb.tokens
		timeNeeded := time.Duration(tokensNeeded) * time.Second / time.Duration(tb.rate)

		// 타임아웃 전에 토큰을 얻을 수 있으면 대기했다가 다시 시도
		if time.Now().Add(timeNeeded).Before(deadline) {
			tb.mutex.Unlock()

			waitTime := minDuration(timeNeeded, deadline.Sub(time.Now()))
			time.Sleep(waitTime)
		} else {
			tb.mutex.Unlock()
			return false
		}

	}
}

func (tb *TokenBucket) TakeWithBurstLimit(tokens, maxBurst int) bool {
	// 최대 버스트보다 많이 가져오지 않도록
	if tokens > maxBurst {
		tokens = maxBurst
	}
	return tb.Take(tokens)
}

func (tb *TokenBucket) Metrics() (capacity, rate, currentTokens int) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()
	return tb.capacity, tb.rate, tb.tokens
}

func (tb *TokenBucket) SetRate(newRate int) {
	tb.mutex.Lock()
	defer tb.mutex.Unlock()

	tb.refill()

	tb.rate = newRate
}

func minDuration(a, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
