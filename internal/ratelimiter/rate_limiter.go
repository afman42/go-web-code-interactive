package ratelimiter

import (
	"net/http"
	"sync"
	"time"
)

// Constants for rate limiter configuration
const (
	CleanupInterval = 10 * time.Minute
	MaxInactivityDuration = 30 * time.Minute
)

// RateLimiter implements a simple rate limiting mechanism using a sliding window counter
type RateLimiter struct {
	visitors map[string][]time.Time
	mu       sync.RWMutex
	window   time.Duration
	maxReqs  int
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(window time.Duration, maxReqs int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string][]time.Time),
		window:   window, // time window (e.g., 1 minute)
		maxReqs:  maxReqs, // max requests per window
	}
	
	// Start cleanup goroutine
	go rl.cleanupVisitors()
	
	return rl
}

// Allow checks if a request from the given IP is allowed
func (rl *RateLimiter) Allow(ip string) bool {
	now := time.Now()
	
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	requests, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = []time.Time{now}
		return true
	}
	
	// Filter out old requests that are outside the window
	updatedRequests := []time.Time{}
	for _, reqTime := range requests {
		if now.Sub(reqTime) < rl.window {
			updatedRequests = append(updatedRequests, reqTime)
		}
	}
	
	// Add current request
	updatedRequests = append(updatedRequests, now)
	rl.visitors[ip] = updatedRequests
	
	// Check if we're within the limit
	return len(updatedRequests) <= rl.maxReqs
}

// cleanupVisitors removes visitors that haven't been seen recently
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(CleanupInterval)
		
		rl.mu.Lock()
		for ip, requests := range rl.visitors {
			// Remove users with no recent requests
			if len(requests) == 0 || time.Since(requests[len(requests)-1]) > MaxInactivityDuration {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// GetVisitorIP extracts the real IP address from the request
func GetVisitorIP(r *http.Request) string {
	// Check X-Forwarded-For header
	forwarded := r.Header.Get("X-Forwarded-For")
	if forwarded != "" {
		// Take the first IP if there are multiple
		for _, ip := range []string{forwarded} {
			return ip
		}
	}

	// Check X-Real-IP header
	realIP := r.Header.Get("X-Real-IP")
	if realIP != "" {
		return realIP
	}

	// Fallback to RemoteAddr
	ip := r.RemoteAddr
	return ip
}

// RateLimitHandler wraps an HTTP handler with rate limiting
func (rl *RateLimiter) RateLimitHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := GetVisitorIP(r)
		if !rl.Allow(ip) {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		
		next.ServeHTTP(w, r)
	})
}