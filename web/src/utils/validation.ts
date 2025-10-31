export function validateUserCode(code: string): {
  isValid: boolean;
  errors: string[];
} {
  const errors: string[] = [];

  // Additional validations
  if (code.length > 10000) {
    // 10KB limit
    errors.push("Code exceeds maximum allowed length (10,000 characters)");
  }

  // Check for basic dangerous patterns that could apply to any language
  if (/\|\|.*;/g.test(code) || /`.*`/g.test(code)) {
    errors.push("Potentially dangerous command execution patterns detected");
  }

  return { isValid: errors.length === 0, errors };
}
