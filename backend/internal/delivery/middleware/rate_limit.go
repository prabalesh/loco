package middleware

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	requests map[string][]time.Time
	mu       sync.Mutex
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{
		requests: make(map[string][]time.Time),
		limit:    limit,
		window:   window,
	}
}

func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get user ID or IP
		userID, ok := GetUserID(r.Context())
		identifier := r.RemoteAddr
		if ok {
			identifier = fmt.Sprintf("user_%d", userID)
		}

		rl.mu.Lock()
		defer rl.mu.Unlock()

		now := time.Now()

		// Clean old requests
		if timestamps, ok := rl.requests[identifier]; ok {
			filtered := []time.Time{}
			for _, t := range timestamps {
				if now.Sub(t) < rl.window {
					filtered = append(filtered, t)
				}
			}
			rl.requests[identifier] = filtered
		}

		// Check limit
		if len(rl.requests[identifier]) >= rl.limit {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		// Add request
		rl.requests[identifier] = append(rl.requests[identifier], now)

		next(w, r)
	}
}
