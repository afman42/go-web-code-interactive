package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
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
}
