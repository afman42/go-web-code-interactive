package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestShellout verifies that code execution for each language works as expected.
func TestShellout(t *testing.T) {
	// Test case for Node.js
	t.Run("Node.js Execution", func(t *testing.T) {
		stdout, stderr, err := Shellout("node", "-e", "console.log('hello node');")
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if strings.TrimSpace(stdout) != "hello node" {
			t.Errorf("Expected stdout 'hello node', but got: %s", stdout)
		}
		if stderr != "" {
			t.Errorf("Expected empty stderr, but got: %s", stderr)
		}
	})

	// Test case for PHP
	t.Run("PHP Execution", func(t *testing.T) {
		stdout, stderr, err := Shellout("php", "-r", "echo 'hello php';")
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if stdout != "hello php" {
			t.Errorf("Expected stdout 'hello php', but got: %s", stdout)
		}
		if stderr != "" {
			t.Errorf("Expected empty stderr, but got: %s", stderr)
		}
	})

	// Test case for Go
	t.Run("Go Execution", func(t *testing.T) {
		// Create a temporary Go file
		tmpFile, err := os.CreateTemp("", "test-*.go")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		goCode := `package main; import "fmt"; func main() { fmt.Println("hello go") }`
		if _, err := tmpFile.WriteString(goCode); err != nil {
			t.Fatalf("Failed to write to temp file: %v", err)
		}
		tmpFile.Close()

		stdout, stderr, err := Shellout("go", "run", tmpFile.Name())
		if err != nil {
			t.Errorf("Expected no error, but got: %v", err)
		}
		if strings.TrimSpace(stdout) != "hello go" {
			t.Errorf("Expected stdout 'hello go', but got: %s", stdout)
		}
		if stderr != "" {
			t.Errorf("Expected empty stderr, but got: %s", stderr)
		}
	})

	// Test for stderr output
	t.Run("Stderr Output", func(t *testing.T) {
		_, stderr, _ := Shellout("node", "-e", "console.error('test error');")
		if !strings.Contains(stderr, "test error") {
			t.Errorf("Expected stderr to contain 'test error', but got: %s", stderr)
		}
	})
}

// TestStringWithCharset ensures the generated string has the correct length.
func TestStringWithCharset(t *testing.T) {
	testCases := []int{5, 10, 20}
	for _, length := range testCases {
		result := StringWithCharset(length)
		if len(result) != length {
			t.Errorf("Expected string of length %d, but got length %d", length, len(result))
		}
	}
}

// TestMoveFile verifies that a file is moved correctly.
func TestMoveFile(t *testing.T) {
	// Create a dummy source file
	source, err := os.CreateTemp("", "source-*.txt")
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	source.WriteString("test content")
	source.Close()

	destPath := filepath.Join(os.TempDir(), "dest-file.txt")
	defer os.Remove(destPath)

	err = MoveFile(source.Name(), destPath)
	if err != nil {
		t.Fatalf("MoveFile failed: %v", err)
	}

	// Check if source file is removed
	if _, err := os.Stat(source.Name()); !os.IsNotExist(err) {
		t.Errorf("Expected source file to be removed, but it still exists.")
	}

	// Check if destination file exists and has content
	content, err := os.ReadFile(destPath)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	if string(content) != "test content" {
		t.Errorf("Expected destination file content 'test content', but got '%s'", string(content))
	}
}

// TestCheckIsNotData validates the slice searching logic.
func TestCheckIsNotData(t *testing.T) {
	slice := []string{"php", "node", "go"}

	if !CheckIsNotData(slice, "node") {
		t.Error("Expected to find 'node' in the slice, but it was not found.")
	}

	if CheckIsNotData(slice, "python") {
		t.Error("Did not expect to find 'python' in the slice, but it was found.")
	}
}
