package security

import (
	"testing"
)

func TestSecurityValidator_ValidateCode(t *testing.T) {
	validator := NewSecurityValidator()

	tests := []struct {
		name     string
		code     string
		language string
		shouldPass bool
	}{
		// Safe JavaScript/Node.js code
		{
			name:     "Safe JavaScript code",
			code:     "console.log('Hello, World!');",
			language: "node",
			shouldPass: true,
		},
		{
			name:     "JavaScript math operations",
			code:     "let a = 1 + 2; console.log(a);",
			language: "node",
			shouldPass: true,
		},
		{
			name:     "JavaScript string operations",
			code:     "let str = 'Hello'; console.log(str.toUpperCase());",
			language: "node",
			shouldPass: true,
		},
		
		// Dangerous JavaScript/Node.js code
		{
			name:     "JavaScript require fs module",
			code:     "const fs = require('fs');",
			language: "node",
			shouldPass: false,
		},
		{
			name:     "JavaScript require child_process module",
			code:     "const exec = require('child_process').exec;",
			language: "node",
			shouldPass: false,
		},
		{
			name:     "JavaScript eval usage",
			code:     "eval('1 + 1');",
			language: "node",
			shouldPass: false,
		},
		{
			name:     "JavaScript process usage",
			code:     "console.log(process.env);",
			language: "node",
			shouldPass: false,
		},
		
		// Safe PHP code
		{
			name:     "Safe PHP code",
			code:     "<?php echo 'Hello, World!'; ?>",
			language: "php",
			shouldPass: true,
		},
		{
			name:     "PHP math operations",
			code:     "<?php $a = 1 + 2; echo $a; ?>",
			language: "php",
			shouldPass: true,
		},
		
		// Dangerous PHP code
		{
			name:     "PHP exec function",
			code:     "<?php exec('ls -la'); ?>",
			language: "php",
			shouldPass: false,
		},
		{
			name:     "PHP system function",
			code:     "<?php system('rm -rf /'); ?>",
			language: "php",
			shouldPass: false,
		},
		{
			name:     "PHP shell_exec function",
			code:     "<?php echo shell_exec('whoami'); ?>",
			language: "php",
			shouldPass: false,
		},
		{
			name:     "PHP file_get_contents with phar",
			code:     "<?php echo file_get_contents('phar://test.phar'); ?>",
			language: "php",
			shouldPass: false,
		},
		
		// Safe Go code
		{
			name:     "Safe Go code",
			code:     "package main\n\nfunc main() {\n    println(\"Hello, World!\")\n}",
			language: "go",
			shouldPass: true,
		},
		{
			name:     "Go basic operations",
			code:     "package main\n\nfunc add(a, b int) int {\n    return a + b\n}\n\nfunc main() {\n    println(add(2, 3))\n}",
			language: "go",
			shouldPass: true,
		},
		
		// Dangerous Go code
		{
			name:     "Go import os package",
			code:     "package main\n\nimport \"os\"\n\nfunc main() {\n    file, _ := os.Open(\"test.txt\")\n}",
			language: "go",
			shouldPass: false,
		},
		{
			name:     "Go import os/exec package",
			code:     "package main\n\nimport \"os/exec\"\n\nfunc main() {\n    cmd := exec.Command(\"ls\")\n    cmd.Run()\n}",
			language: "go",
			shouldPass: false,
		},
		{
			name:     "Go import net/http package",
			code:     "package main\n\nimport \"net/http\"\n\nfunc main() {\n    http.Get(\"http://example.com\")\n}",
			language: "go",
			shouldPass: false,
		},
		
		// Common dangerous patterns
		{
			name:     "Directory traversal",
			code:     "../../../etc/passwd",
			language: "node",
			shouldPass: false,
		},
		{
			name:     "System path access",
			code:     "/etc/passwd",
			language: "node",
			shouldPass: false,
		},
		
		// Unsupported language
		{
			name:     "Unsupported language",
			code:     "print('Hello')",
			language: "python",
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCode(tt.code, tt.language)
			if tt.shouldPass && err != nil {
				t.Errorf("Expected validation to pass, but got error: %v", err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Expected validation to fail, but got no error")
			}
		})
	}
}

func TestSecurityValidator_ValidateCodeEmptyCode(t *testing.T) {
	validator := NewSecurityValidator()
	
	err := validator.ValidateCode("", "node")
	if err != nil {
		t.Errorf("Expected empty code to pass validation, but got error: %v", err)
	}
}

func TestSecurityValidator_ValidateCodeCaseInsensitive(t *testing.T) {
	validator := NewSecurityValidator()
	
	// Test that validation works case-insensitively
	tests := []struct {
		name     string
		code     string
		language string
		shouldPass bool
	}{
		{
			name:     "Case insensitive exec detection",
			code:     "EXEC('ls');",
			language: "node",
			shouldPass: false,
		},
		{
			name:     "Case insensitive PHP system",
			code:     "<?php SYSTEM('rm -rf /'); ?>",
			language: "php",
			shouldPass: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.ValidateCode(tt.code, tt.language)
			if tt.shouldPass && err != nil {
				t.Errorf("Expected validation to pass, but got error: %v", err)
			}
			if !tt.shouldPass && err == nil {
				t.Errorf("Expected validation to fail, but got no error")
			}
		})
	}
}