package ratelimiter

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRateLimiter_Allow(t *testing.T) {
	// Use different IPs to avoid conflicts between tests
	rl := NewRateLimiter(1*time.Second, 2) // Allow 2 requests per second
	
	// Should allow first two requests
	if !rl.Allow("127.0.0.1") {
		t.Error("Expected first request to be allowed")
	}
	if !rl.Allow("127.0.0.1") {
		t.Error("Expected second request to be allowed")
	}
	
	// Should block the third request immediately since we've already made 2 in the window
	if rl.Allow("127.0.0.1") {
		t.Error("Expected third request to be blocked")
	}
	
	// Test with a different IP to avoid conflicts from the previous test
	if !rl.Allow("127.0.0.2") {
		t.Error("Expected request from different IP to be allowed")
	}
}

func TestRateLimitHandler(t *testing.T) {
	rl := NewRateLimiter(100*time.Millisecond, 1) // Allow 1 request per 100 milliseconds
	
	// Create a test handler that always succeeds
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	})
	
	limitedHandler := rl.RateLimitHandler(nextHandler)
	
	// First request should succeed
	req1 := httptest.NewRequest("GET", "/", nil)
	req1.RemoteAddr = "127.0.0.1:12345" // Set a specific remote address for testing
	w1 := httptest.NewRecorder()
	limitedHandler.ServeHTTP(w1, req1)
	
	if w1.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w1.Code)
	}
	
	// Second request should be rate limited (sent immediately after first)
	req2 := httptest.NewRequest("GET", "/", nil)
	req2.RemoteAddr = "127.0.0.1:12345" // Same IP as first request
	w2 := httptest.NewRecorder()
	limitedHandler.ServeHTTP(w2, req2)
	
	if w2.Code != http.StatusTooManyRequests {
		t.Errorf("Expected status %d, got %d", http.StatusTooManyRequests, w2.Code)
	}
}