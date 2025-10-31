export function validateUserCode(code: string, language: string): { isValid: boolean; errors: string[] } {
  // Static patterns that won't change
  const dangerousPatterns = {
    node: [
      /require\(['"`]\s*child_process\s*['"`]\)/,
      /require\(['"`]\s*fs\s*['"`]\)/,
      /exec\s*\(/,
      /eval\s*\(/,
      /setTimeout\s*\(\s*['"`].*['"`]/,
      /setInterval\s*\(\s*['"`].*['"`]/,
    ],
    php: [
      /exec\s*\(/,
      /system\s*\(/,
      /shell_exec\s*\(/,
      /passthru\s*\(/,
      /proc_open\s*\(/,
      /file_get_contents\s*\(['"`]phar:/,
      /file_get_contents\s*\(['"`]ssh2:/,
      /file_put_contents\s*\(['"`]php:\/\/.*['"`]\)/,
      /fopen\s*\(['"`]php:\/\/.*['"`]\)/,
    ],
    go: [
      /import.*"os"/,
      /import.*"os\/exec"/,
      /import.*"net"/,
      /import.*"net\/http"/,
      /import.*"syscall"/,
      /os\.Open\s*\(/,
      /os\.Create\s*\(/,
      /exec\.Command\s*\(/,
      /http\.Get\s*\(/,
    ]
  };
  
  const patterns = dangerousPatterns[language as keyof typeof dangerousPatterns] || [];
  const errors: string[] = [];
  
  for (const pattern of patterns) {
    if (pattern.test(code)) {
      errors.push(`Potentially dangerous pattern detected: ${pattern.toString()}`);
    }
  }
  
  // Additional validations
  if (code.length > 10000) { // 10KB limit
    errors.push('Code exceeds maximum allowed length (10,000 characters)');
  }
  
  // Check for basic dangerous patterns that could apply to any language
  if (/\|\|.*;/g.test(code) || /`.*`/g.test(code)) {
    errors.push('Potentially dangerous command execution patterns detected');
  }
  
  return { isValid: errors.length === 0, errors };
}