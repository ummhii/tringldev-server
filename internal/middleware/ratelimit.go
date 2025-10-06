package middleware

import (
	"log"
	"sync"
	"time"

	"github.com/kataras/iris/v12"
)

type RateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     time.Duration
	burst    int
}

type visitor struct {
	limiter  *tokenBucket
	lastSeen time.Time
}

type tokenBucket struct {
	tokens    float64
	maxTokens int
	refillAt  time.Time
	rate      time.Duration
	mu        sync.Mutex
}

func NewRateLimiter(rate time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
	}

	go rl.cleanupVisitors()

	return rl
}

func (tb *tokenBucket) allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	now := time.Now()

	if now.After(tb.refillAt) {
		elapsed := now.Sub(tb.refillAt)
		tokensToAdd := float64(elapsed) / float64(tb.rate)
		tb.tokens = min(tb.tokens+tokensToAdd, float64(tb.maxTokens))
		tb.refillAt = now
	}

	if tb.tokens >= 1 {
		tb.tokens--
		return true
	}

	return false
}

func (rl *RateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			limiter: &tokenBucket{
				tokens:    float64(rl.burst),
				maxTokens: rl.burst,
				refillAt:  time.Now(),
				rate:      rl.rate,
			},
			lastSeen: time.Now(),
		}
		rl.visitors[ip] = v
	}

	v.lastSeen = time.Now()
	return v
}

func (rl *RateLimiter) cleanupVisitors() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Handler() iris.Handler {
	return func(ctx iris.Context) {
		ip := ctx.RemoteAddr()

		trustedIPs := map[string]bool{
			"127.0.0.1": true, // Localhost
		}

		if trustedIPs[ip] {
			ctx.Next()
			return
		}

		v := rl.getVisitor(ip)

		if !v.limiter.allow() {
			log.Printf("Rate limit exceeded for IP: %s", ip)
			ctx.StatusCode(iris.StatusTooManyRequests)
			ctx.JSON(iris.Map{
				"error": "Rate limit exceeded. Please try again later.",
			})
			return
		}

		ctx.Next()
	}
}

func Stricter(rate time.Duration, burst int) *RateLimiter {
	return NewRateLimiter(rate, burst)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
