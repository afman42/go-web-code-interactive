package security

import (
	"fmt"
	"regexp"
	"strings"
)

// SecurityValidator validates user code for potentially dangerous operations
type SecurityValidator struct {
	forbiddenPatterns map[string][]*regexp.Regexp
	globalPatterns    []*regexp.Regexp
	allowedLanguages  []string
}

// NewSecurityValidator creates a new security validator
func NewSecurityValidator() *SecurityValidator {
	// Language-specific forbidden patterns
	forbiddenPatterns := make(map[string][]*regexp.Regexp)
	
	// JavaScript/Node.js patterns
	forbiddenPatterns["node"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(require\(\s*["']child_process["']\s*\))`), // child_process
		regexp.MustCompile(`(?i)(require\(\s*["']fs["']\s*\))`), // fs module
		regexp.MustCompile(`(?i)(require\(\s*["']os["']\s*\))`), // os module
		regexp.MustCompile(`(?i)(require\(\s*["']net["']\s*\))`), // net module
		regexp.MustCompile(`(?i)(require\(\s*["']dns["']\s*\))`), // dns module
		regexp.MustCompile(`(?i)(require\(\s*["']http["']\s*\))`), // http module
		regexp.MustCompile(`(?i)(require\(\s*["']https["']\s*\))`), // https module
		regexp.MustCompile(`(?i)(eval\s*\()`), // eval
		regexp.MustCompile(`(?i)(EXEC\s*\()`), // EXEC function
		regexp.MustCompile(`(?i)(Function\s*\(\s*["'].*["']\s*\))`), // Function constructor
		regexp.MustCompile(`(?i)(setTimeout\s*\(\s*["'].*["']\s*,?)`), // setTimeout with string
		regexp.MustCompile(`(?i)(setInterval\s*\(\s*["'].*["']\s*,?)`), // setInterval with string
		regexp.MustCompile(`(?i)(import\(\s*["'].+["']\s*\))`), // dynamic imports
		regexp.MustCompile(`(?i)(process\.)`), // process object usage
		regexp.MustCompile(`(?i)(global\.)`), // global object usage
		regexp.MustCompile(`(?i)(__dirname|__filename)`), // file system paths
		regexp.MustCompile(`(?i)(new\s+ActiveXObject\s*\()`), // ActiveXObject (if somehow available)
		regexp.MustCompile(`(?i)(WScript\.)`), // WScript (if somehow available)
	}
	
	// PHP patterns
	forbiddenPatterns["php"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(exec\s*\()`), // exec
		regexp.MustCompile(`(?i)(shell_exec\s*\()`), // shell_exec
		regexp.MustCompile(`(?i)(system\s*\()`), // system
		regexp.MustCompile(`(?i)(passthru\s*\()`), // passthru
		regexp.MustCompile(`(?i)(popen\s*\()`), // popen
		regexp.MustCompile(`(?i)(proc_open\s*\()`), // proc_open
		regexp.MustCompile(`(?i)(pcntl_exec\s*\()`), // pcntl_exec
		regexp.MustCompile(`(?i)(pcntl_fork\s*\()`), // pcntl_fork
		regexp.MustCompile(`(?i)(dl\s*\()`), // dl
		regexp.MustCompile(`(?i)(file_get_contents\s*\(\s*['"]phar:)`), // phar protocol
		regexp.MustCompile(`(?i)(file_get_contents\s*\(\s*['"]ssh2:)`), // ssh2 protocol
		regexp.MustCompile(`(?i)(file_get_contents\s*\(\s*['"]compress\.zlib:)`), // compress.zlib protocol
		regexp.MustCompile(`(?i)(file_get_contents\s*\(\s*['"]compress\.bzip2:)`), // compress.bzip2 protocol
		regexp.MustCompile(`(?i)(file_get_contents\s*\(\s*['"]glob:)`), // glob protocol
		regexp.MustCompile(`(?i)(file_get_contents\s*\(\s*['"]data:)`), // data protocol
		regexp.MustCompile(`(?i)(file_put_contents\s*\(\s*[\'"]php://.*[\'"])`), // php:// wrappers
		regexp.MustCompile(`(?i)(fopen\s*\(\s*[\'"]php://.*[\'"])`), // php:// wrappers
		regexp.MustCompile(`(?i)(include\s+[\'"]phar:)`), // phar in include
		regexp.MustCompile(`(?i)(include_once\s+[\'"]phar:)`), // phar in include_once
		regexp.MustCompile(`(?i)(require\s+[\'"]phar:)`), // phar in require
		regexp.MustCompile(`(?i)(require_once\s+[\'"]phar:)`), // phar in require_once
		regexp.MustCompile(`(?i)(scandir\s*\()`), // scandir
		regexp.MustCompile(`(?i)(opendir\s*\()`), // opendir
		regexp.MustCompile(`(?i)(glob\s*\()`), // glob
		regexp.MustCompile(`(?i)(exec\s*\()`), // exec
		regexp.MustCompile(`(?i)(passthru\s*\()`), // passthru
		regexp.MustCompile(`(?i)($_(GET|POST|COOKIE|SESSION|SERVER|FILES|ENV))`), // superglobals
		regexp.MustCompile(`(?i)(getenv\s*\()`), // getenv
		regexp.MustCompile(`(?i)(putenv\s*\()`), // putenv
		regexp.MustCompile(`(?i)(apache_setenv\s*\()`), // apache_setenv
		regexp.MustCompile(`(?i)(apache_getenv\s*\()`), // apache_getenv
	}
	
	// Go patterns
	forbiddenPatterns["go"] = []*regexp.Regexp{
		regexp.MustCompile(`(?i)(import\s+\(\s*.*"os".*\))`), // os import
		regexp.MustCompile(`(?i)(import\s+.*"os"\s*)`), // os import
		regexp.MustCompile(`(?i)(import\s+\(\s*.*"os/exec".*\))`), // os/exec import
		regexp.MustCompile(`(?i)(import\s+.*"os/exec"\s*)`), // os/exec import
		regexp.MustCompile(`(?i)(import\s+\(\s*.*"net".*\))`), // net import
		regexp.MustCompile(`(?i)(import\s+.*"net"\s*)`), // net import
		regexp.MustCompile(`(?i)(import\s+\(\s*.*"net/http".*\))`), // net/http import
		regexp.MustCompile(`(?i)(import\s+.*"net/http"\s*)`), // net/http import
		regexp.MustCompile(`(?i)(import\s+\(\s*.*"syscall".*\))`), // syscall import
		regexp.MustCompile(`(?i)(import\s+.*"syscall"\s*)`), // syscall import
		regexp.MustCompile(`(?i)(import\s+\(\s*.*"unsafe".*\))`), // unsafe import
		regexp.MustCompile(`(?i)(import\s+.*"unsafe"\s*)`), // unsafe import
		regexp.MustCompile(`(?i)(import\s+\(\s*.*"os/user".*\))`), // os/user import
		regexp.MustCompile(`(?i)(import\s+.*"os/user"\s*)`), // os/user import
		regexp.MustCompile(`(?i)(os\.Open\s*\()`), // os.Open
		regexp.MustCompile(`(?i)(os\.Create\s*\()`), // os.Create
		regexp.MustCompile(`(?i)(os\.Remove\s*\()`), // os.Remove
		regexp.MustCompile(`(?i)(os\.Mkdir\s*\()`), // os.Mkdir
		regexp.MustCompile(`(?i)(os\.MkdirAll\s*\()`), // os.MkdirAll
		regexp.MustCompile(`(?i)(os\.Rename\s*\()`), // os.Rename
		regexp.MustCompile(`(?i)(os\.Chmod\s*\()`), // os.Chmod
		regexp.MustCompile(`(?i)(os\.Chown\s*\()`), // os.Chown
		regexp.MustCompile(`(?i)(os\.Symlink\s*\()`), // os.Symlink
		regexp.MustCompile(`(?i)(exec\.Command\s*\()`), // exec.Command
		regexp.MustCompile(`(?i)(http\.NewRequest\s*\()`), // http.NewRequest
		regexp.MustCompile(`(?i)(http\.Get\s*\()`), // http.Get
		regexp.MustCompile(`(?i)(http\.Post\s*\()`), // http.Post
		regexp.MustCompile(`(?i)(http\.Do\s*\()`), // http.Do
		regexp.MustCompile(`(?i)(syscall\.)`), // syscall usage
		regexp.MustCompile(`(?i)(unsafe\.)`), // unsafe usage
		regexp.MustCompile(`(?i)(user\.Current\s*\()`), // user.Current
	}
	
	// Global patterns that apply to all languages
	globalPatterns := []*regexp.Regexp{
		// Directory traversal
		regexp.MustCompile(`(?i)(\.\.\/|\.\.\\)`),
		
		// Sensitive system paths
		regexp.MustCompile(`(?i)(\/etc\/|\/proc\/|\/sys\/|\/dev\/|C:\\Windows\\|C:\\ProgramData\\|C:\\Program Files)`),
	}

	return &SecurityValidator{
		forbiddenPatterns: forbiddenPatterns,
		globalPatterns:    globalPatterns,
		allowedLanguages:  []string{"node", "php", "go"},
	}
}

// ValidateCode validates the user code for security issues
func (sv *SecurityValidator) ValidateCode(code string, language string) error {
	// Check if language is supported
	isValidLang := false
	for _, lang := range sv.allowedLanguages {
		if strings.ToLower(lang) == strings.ToLower(language) {
			isValidLang = true
			break
		}
	}
	
	if !isValidLang {
		return fmt.Errorf("unsupported language: %s", language)
	}

	// Check language-specific patterns
	if patterns, exists := sv.forbiddenPatterns[strings.ToLower(language)]; exists {
		for _, pattern := range patterns {
			if pattern.MatchString(code) {
				return fmt.Errorf("forbidden operation detected in %s: %s", language, pattern.String())
			}
		}
	}

	// Check global patterns
	for _, pattern := range sv.globalPatterns {
		if pattern.MatchString(code) {
			return fmt.Errorf("forbidden operation detected: %s", pattern.String())
		}
	}

	return nil
}