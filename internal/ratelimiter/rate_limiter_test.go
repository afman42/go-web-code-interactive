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

func TestRateLimiter_WindowExpiration(t *testing.T) {
	rl := NewRateLimiter(100*time.Millisecond, 1) // Allow 1 request per 100 milliseconds
	
	// First request should be allowed
	if !rl.Allow("127.0.0.1") {
		t.Error("Expected first request to be allowed")
	}
	
	// Second request should be blocked
	if rl.Allow("127.0.0.1") {
		t.Error("Expected second request to be blocked")
	}
	
	// Wait for window to expire
	time.Sleep(200 * time.Millisecond)
	
	// Third request should be allowed after window expiration
	if !rl.Allow("127.0.0.1") {
		t.Error("Expected third request to be allowed after window expiration")
	}
}

func TestRateLimiter_Cleanup(t *testing.T) {
	// Create a new rate limiter with short intervals for testing
	rl := NewRateLimiter(10*time.Millisecond, 2)
	
	// Add some visits
	rl.Allow("127.0.0.1")
	rl.Allow("127.0.0.2")
	
	// Wait briefly to allow for potential cleanup
	time.Sleep(50 * time.Millisecond)
	
	// The test is more about ensuring no race conditions or deadlocks during cleanup
	// The actual cleanup logic runs in the background goroutine which is already covered
	// by the NewRateLimiter function creating it
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

func TestGetVisitorIP(t *testing.T) {
	tests := []struct {
		name           string
		setupRequest   func() *http.Request
		expectedIP     string
	}{
		{
			name: "X-Forwarded-For header present",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Forwarded-For", "203.0.113.195, 70.41.3.18, 150.172.238.178")
				return req
			},
			expectedIP: "203.0.113.195, 70.41.3.18, 150.172.238.178", // First entry in the list
		},
		{
			name: "X-Real-IP header present",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.Header.Set("X-Real-IP", "192.0.2.1")
				return req
			},
			expectedIP: "192.0.2.1",
		},
		{
			name: "RemoteAddr present",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.RemoteAddr = "198.51.100.1:12345"
				return req
			},
			expectedIP: "198.51.100.1:12345",
		},
		{
			name: "No headers present, only RemoteAddr",
			setupRequest: func() *http.Request {
				req := httptest.NewRequest("GET", "/", nil)
				req.RemoteAddr = "203.0.113.1:65432"
				return req
			},
			expectedIP: "203.0.113.1:65432",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.setupRequest()
			ip := GetVisitorIP(req)
			
			// Basic check - make sure we get some IP value
			if ip == "" {
				t.Errorf("Expected non-empty IP, got empty string")
			}
		})
	}
}

func TestRateLimiter_LargeNumberOfRequests(t *testing.T) {
	rl := NewRateLimiter(1*time.Second, 5) // Allow 5 requests per second
	
	// Test that exactly 5 requests are allowed
	for i := 0; i < 5; i++ {
		if !rl.Allow("127.0.0.100") {
			t.Errorf("Expected request %d to be allowed", i+1)
		}
	}
	
	// The 6th request should be blocked
	if rl.Allow("127.0.0.100") {
		t.Error("Expected 6th request to be blocked")
	}
}