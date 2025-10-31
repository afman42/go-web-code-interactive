package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

// Setup: Set environment variables required for the tests.
func TestMain(m *testing.M) {
	os.Setenv("CORS_DOMAIN", "http://localhost:5173")
	os.Setenv("APP_PORT", "8000")
	// Create tmp directory if it doesn't exist
	if _, err := os.Stat("./tmp"); os.IsNotExist(err) {
		os.Mkdir("./tmp", os.ModePerm)
	}
	code := m.Run()
	os.Exit(code)
}

// TestIndexHandler covers all logic within the main router.
func TestIndexHandler(t *testing.T) {
	// Test GET request
	t.Run("GET Request", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}
		
		// Check if response contains HTML content
		if !strings.Contains(rr.Header().Get("Content-Type"), "text/html") {
			t.Errorf("Expected Content-Type to contain text/html, got: %s", rr.Header().Get("Content-Type"))
		}
	})

	// Test successful POST request
	t.Run("POST Request - Success", func(t *testing.T) {
		jsonData := `{"txt": "console.log(\"hello world\");", "lang": "node", "type": "repl"}`
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
		}

		expectedBody := `"out":"hello world\n"`
		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedBody)) {
			t.Errorf("handler returned unexpected body: got %v want it to contain %v", rr.Body.String(), expectedBody)
		}
	})

	// Test POST request with STQ mode
	t.Run("POST Request - STQ Success", func(t *testing.T) {
		jsonData := `{"txt": "function intIntoString(n) { return n.toString(); }", "lang": "node", "type": "stq"}`
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		// Should return 200 for valid STQ code (though output might be different)
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code for STQ: got %v want %v", status, http.StatusOK)
		}
	})

	// Test POST request with invalid language
	t.Run("POST Request - Invalid Language", func(t *testing.T) {
		jsonData := `{"txt": "print(\"hello\")", "lang": "python", "type": "repl"}`
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for invalid lang: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test POST request with invalid type
	t.Run("POST Request - Invalid Type", func(t *testing.T) {
		jsonData := `{"txt": "console.log(\"hello\")", "lang": "node", "type": "invalid"}`
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for invalid type: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test POST request with empty fields
	t.Run("POST Request - Empty Field", func(t *testing.T) {
		jsonData := `{"txt": "", "lang": "node", "type": "repl"}` // Empty txt field
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for empty field: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test POST request with missing fields
	t.Run("POST Request - Missing Fields", func(t *testing.T) {
		jsonData := `{"txt": "console.log(\"hello\")", "lang": "node"}` // Missing type field
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for missing field: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test POST request with malformed JSON
	t.Run("POST Request - Malformed JSON", func(t *testing.T) {
		jsonData := `{"txt": "console.log("hello")", "lang": "node", "type": "repl"}` // Invalid JSON - unescaped quotes
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for malformed JSON: got %v want %v", status, http.StatusBadRequest)
		}
	})

	// Test POST request with security validation failure
	t.Run("POST Request - Security Validation Failure", func(t *testing.T) {
		jsonData := `{"txt": "const fs = require('fs'); fs.readFileSync('/etc/passwd');", "lang": "node", "type": "repl"}`
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusBadRequest {
			t.Errorf("handler returned wrong status code for security violation: got %v want %v", status, http.StatusBadRequest)
		}

		// Check if response body contains security validation error
		if !strings.Contains(rr.Body.String(), "Security validation failed") {
			t.Errorf("expected response to contain 'Security validation failed', got: %s", rr.Body.String())
		}
	})

	// Test unsupported HTTP method
	t.Run("Unsupported Method", func(t *testing.T) {
		req, err := http.NewRequest("PUT", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusMethodNotAllowed {
			t.Errorf("handler returned wrong status code for unsupported method: got %v want %v", status, http.StatusMethodNotAllowed)
		}
	})

	// Test OPTIONS request
	t.Run("OPTIONS Request", func(t *testing.T) {
		req, err := http.NewRequest("OPTIONS", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNoContent {
			t.Errorf("handler returned wrong status code for OPTIONS: got %v want %v", status, http.StatusNoContent)
		}
		
		// Check that essential CORS headers are set for OPTIONS requests
		// The Access-Control-Allow-Origin header is set to IPCors which is loaded from the environment
		// In test environment, we just make sure the essential headers are set
		if rr.Header().Get("Access-Control-Allow-Methods") == "" {
			t.Error("Expected Access-Control-Allow-Methods header to be set")
		}
		if rr.Header().Get("Access-Control-Allow-Headers") == "" {
			t.Error("Expected Access-Control-Allow-Headers header to be set")
		}
	})

	// Test non-root path (should return 404)
	t.Run("Non-Root Path", func(t *testing.T) {
		req, err := http.NewRequest("GET", "/not-found", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("handler returned wrong status code for non-root path: got %v want %v", status, http.StatusNotFound)
		}
	})

	// Test PHP execution
	t.Run("POST Request - PHP Success", func(t *testing.T) {
		jsonData := `{"txt": "<?php echo 'hello php'; ?>", "lang": "php", "type": "repl"}`
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code for PHP: got %v want %v", status, http.StatusOK)
		}

		expectedBody := `"out":"hello php"`
		if !bytes.Contains(rr.Body.Bytes(), []byte(expectedBody)) && 
			!bytes.Contains(rr.Body.Bytes(), []byte(`"out":"hello php\n"`)) {
			t.Errorf("handler returned unexpected body for PHP: got %v want it to contain %v", rr.Body.String(), expectedBody)
		}
	})

	// Test Go execution
	t.Run("POST Request - Go Success", func(t *testing.T) {
		jsonData := `{"txt": "package main\n\nfunc main() {\n\tprintln(\"hello go\")\n}", "lang": "go", "type": "repl"}`
		req, err := http.NewRequest("POST", "/", bytes.NewBufferString(jsonData))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(index)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusOK {
			t.Errorf("handler returned wrong status code for Go: got %v want %v", status, http.StatusOK)
		}

		// For Go, println() outputs to stderr, so we expect to see "hello go" in errout field
		if !bytes.Contains(rr.Body.Bytes(), []byte(`"errout":"hello go`)) {
			t.Errorf("handler returned unexpected body for Go: got %v", rr.Body.String())
		}
	})

	// The rate limiting test is best handled in the ratelimiter package tests directly
	// since it's already well tested there, and changing the global rateLimiter variable
	// can affect other tests. This test has been moved to the ratelimiter package.
}

// Test utility functions that are used in main.go
func TestDataStruct(t *testing.T) {
	// Test that the Data struct can be marshaled/unmarshaled correctly
	data := Data{
		Txt:        "test code",
		Stdout:     "output",
		Stderr:     "error",
		StatusCode: 200,
		Language:   "node",
		Tipe:       "repl",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal Data struct: %v", err)
	}

	var unmarshaledData Data
	err = json.Unmarshal(jsonData, &unmarshaledData)
	if err != nil {
		t.Fatalf("Failed to unmarshal Data struct: %v", err)
	}

	if unmarshaledData.Txt != data.Txt || unmarshaledData.Language != data.Language {
		t.Errorf("Data struct marshaling/unmarshaling failed: got %+v, want %+v", unmarshaledData, data)
	}
}
