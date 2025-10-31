// Editor constants
export const EDITOR_DEFAULT_HEIGHT = '96vh'; // Use viewport height instead of fixed pixels
export const EDITOR_MIN_HEIGHT = '400px';
export const EDITOR_MAX_HEIGHT = '90vh';
export const CODE_MAX_LENGTH = 10000; // Maximum allowed code length

// API constants
export const API_ENDPOINT_EXECUTE = '/';
export const API_TIMEOUT = 30000; // 30 seconds

// UI constants
export const TOAST_DURATION_SHORT = 1000;
export const TOAST_DURATION_MEDIUM = 3000;
export const TOAST_DURATION_LONG = 5000;

// Language constants
export const SUPPORTED_LANGUAGES = {
  NODE: 'node',
  PHP: 'php',
  GO: 'go'
} as const;

export const EXECUTION_TYPES = {
  REPL: 'repl',
  STQ: 'stq'
} as const;

// Theme constants
export const THEMES = {
  LIGHT: 'light',
  DARK: 'dark'
} as const;

// Severity levels for diagnostics
export const SEVERITY = {
  ERROR: 'error',
  WARNING: 'warning',
  INFO: 'info',
  HINT: 'hint'
} as const;