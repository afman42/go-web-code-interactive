import { CODE_MAX_LENGTH } from '../constants';

export function validateUserCode(code: string): {
  isValid: boolean;
  errors: string[];
} {
  const errors: string[] = [];

  // Additional validations
  if (code.length > CODE_MAX_LENGTH) {
    errors.push(`Code exceeds maximum allowed length (${CODE_MAX_LENGTH} characters)`);
  }

  // Check for basic dangerous patterns that could apply to any language
  if (/\|\|.*;/g.test(code) || /`.*`/g.test(code)) {
    errors.push("Potentially dangerous command execution patterns detected");
  }

  return { isValid: errors.length === 0, errors };
}
