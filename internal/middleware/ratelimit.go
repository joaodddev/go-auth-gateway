package middleware

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	redisClient *redis.Client
	requests    int
	duration    time.Duration
}

func NewRateLimiter(redisClient *redis.Client, requests int, duration time.Duration) *RateLimiter {
	return &RateLimiter{
		redisClient: redisClient,
		requests:    requests,
		duration:    duration,
	}
}

func (rl *RateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr
		ctx := context.Background()
		key := fmt.Sprintf("ratelimit:%s", ip)

		count, err := rl.redisClient.Incr(ctx, key).Result()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if count == 1 {
			rl.redisClient.Expire(ctx, key, rl.duration)
		}

		if count > int64(rl.requests) {
			w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", time.Now().Add(rl.duration).Format(time.RFC1123))
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		w.Header().Set("X-RateLimit-Limit", fmt.Sprintf("%d", rl.requests))
		w.Header().Set("X-RateLimit-Remaining", fmt.Sprintf("%d", rl.requests-int(count)))
		next(w, r)
	}
}
