<script lang="ts">
import { onMount, onDestroy } from "svelte";
import { EditorState, type Extension } from "@codemirror/state";
import {
  EditorView,
  keymap,
  highlightSpecialChars,
  drawSelection,
  highlightActiveLine,
  dropCursor,
  rectangularSelection,
  crosshairCursor,
  lineNumbers,
  highlightActiveLineGutter,
  ViewPlugin,
  ViewUpdate
} from "@codemirror/view";
import {
  defaultHighlightStyle,
  syntaxHighlighting,
  indentOnInput,
  bracketMatching,
  foldGutter,
  foldKeymap,
  type LanguageSupport,
} from "@codemirror/language";
import {
  defaultKeymap,
  history,
  historyKeymap,
  indentWithTab
} from "@codemirror/commands";
import {
  searchKeymap,
  highlightSelectionMatches,
  search
} from "@codemirror/search";
import {
  autocompletion,
  completionKeymap,
  closeBrackets,
  closeBracketsKeymap,
  snippetCompletion,
  type CompletionContext,
  type CompletionResult,
} from "@codemirror/autocomplete";
import { lintKeymap, linter, type Diagnostic } from "@codemirror/lint";
import { javascript } from "@codemirror/lang-javascript";
import { go } from "@codemirror/lang-go";
import { php } from "@codemirror/lang-php";
import { oneDark } from "@codemirror/theme-one-dark";

// Props using Svelte 5 runes syntax
let { 
  language = "node", 
  code = "", 
  onChange = (() => {}) as (value: string) => void, 
  theme = 'light' 
}: {
  language: string;
  code: string;
  onChange: (value: string) => void;
  theme: 'light' | 'dark';
} = $props();

// Reactive states
let view: EditorView | null = null;
let editorContainer: HTMLDivElement;

// Theme extension based on preference
const themeExtension = theme === 'dark' ? [oneDark] : [];

// Language configurations
const languageConfigs: Record<string, LanguageSupport> = {
  node: javascript({ jsx: true, typescript: false }),
  php: php(),
  go: go()
};

// VS Code-like autocompletion for PHP
function phpCompletions(context: CompletionContext): CompletionResult | null {
  const word = context.matchBefore(/\w*/);
  if (!word || word.from === word.to && !context.explicit) return null;

  const phpKeywords = [
    "echo", "else", "elseif", "foreach", "function", "return", "if", "while", "for", "break", "continue",
    "class", "public", "private", "protected", "namespace", "use", "extends", "implements", 
    "interface", "trait", "abstract", "final", "static", "const", "global", "include", "require"
  ];
  const phpFunctions = [
    "array", "strlen", "str_replace", "explode", "implode", "isset", "empty", "array_push",
    "array_pop", "count", "array_keys", "array_values", "array_map", "array_filter",
    "var_dump", "print_r", "json_encode", "json_decode", "file_get_contents", "file_put_contents"
  ];

  return {
    from: word.from,
    options: [
      ...phpKeywords.map((kw: string) => ({ label: kw, type: "keyword", detail: "PHP keyword" })),
      ...phpFunctions.map((fn: string) => ({ label: fn, type: "function", detail: "PHP function" })),
      snippetCompletion("for ($${i} = 0; $${i} < $${len}; $${i}++) {\n\t$${}\n}", { label: "for", detail: "for loop", type: "snippet" }),
      snippetCompletion("foreach ($${array} as $${key} => $${value}) {\n\t$${}\n}", { label: "foreach", detail: "foreach loop", type: "snippet" }),
      snippetCompletion("if ($${condition}) {\n\t$${}\n}", { label: "if", detail: "if statement", type: "snippet" }),
      snippetCompletion("function $${name}($${params}) {\n\t$${}\n}", { label: "function", detail: "function definition", type: "snippet" })
    ]
  };
}

// VS Code-like autocompletion for Go
function goCompletions(context: CompletionContext): CompletionResult | null {
  const word = context.matchBefore(/\w*/);
  if (!word || word.from === word.to && !context.explicit) return null;

  const goKeywords = [
    "func", "var", "const", "else", "for", "range", "return", "if", "switch", "case", 
    "break", "continue", "default", "defer", "go", "goto", "interface", "map", 
    "package", "import", "type", "struct", "chan", "select", "fallthrough"
  ];
  const goBuiltins = [
    "println", "print", "len", "cap", "make", "new", "append", "copy", "close",
    "delete", "complex", "real", "imag", "panic", "recover", "error", "bool", "string"
  ];

  return {
    from: word.from,
    options: [
      ...goKeywords.map((kw: string) => ({ label: kw, type: "keyword", detail: "Go keyword" })),
      ...goBuiltins.map((fn: string) => ({ label: fn, type: "function", detail: "Go builtin" })),
      snippetCompletion("if $${condition} {\n\t$${}\n}", { label: "if", detail: "if statement", type: "snippet" }),
      snippetCompletion("for $${i} := 0; $${i} < $${len}; $${i}++ {\n\t$${}\n}", { label: "for", detail: "for loop", type: "snippet" }),
      snippetCompletion("func $${name}($${params}) $${returnType} {\n\t$${}\n}", { label: "func", detail: "function definition", type: "snippet" }),
      snippetCompletion("switch $${value} {\ncase $${condition}:\n\t$${}\n}", { label: "switch", detail: "switch statement", type: "snippet" })
    ]
  };
}

// JavaScript/Node.js completions
function jsCompletions(context: CompletionContext): CompletionResult | null {
  const word = context.matchBefore(/\w*/);
  if (!word || word.from === word.to && !context.explicit) return null;

  const jsKeywords = [
    "function", "return", "if", "else", "for", "while", "break", "continue", "var", "let", 
    "const", "class", "extends", "super", "new", "this", "import", "export", "default", 
    "async", "await", "try", "catch", "finally", "throw", "switch", "case", "default", 
    "do", "typeof", "instanceof", "in", "of", "with", "debugger", "yield"
  ];
  const jsFunctions = [
    "console", "log", "error", "warn", "info", "debug", "time", "timeEnd", 
    "setTimeout", "setInterval", "clearTimeout", "clearInterval",
    "Array", "Object", "String", "Number", "Boolean", "Date", "Math", "JSON"
  ];

  return {
    from: word.from,
    options: [
      ...jsKeywords.map((kw: string) => ({ label: kw, type: "keyword", detail: "JavaScript keyword" })),
      ...jsFunctions.map((fn: string) => ({ label: fn, type: "function", detail: "JavaScript function" })),
      snippetCompletion("function $${name}($${params}) {\n\t$${}\n}", { label: "function", detail: "function declaration", type: "snippet" }),
      snippetCompletion("($${params}) => {\n\t$${}\n}", { label: "arrow", detail: "arrow function", type: "snippet" }),
      snippetCompletion("for (let $${i} = 0; $${i} < $${len}; $${i}++) {\n\t$${}\n}", { label: "for", detail: "for loop", type: "snippet" }),
      snippetCompletion("if ($${condition}) {\n\t$${}\n}", { label: "if", detail: "if statement", type: "snippet" }),
      snippetCompletion("try {\n\t$${}\n} catch ($${error}) {\n\t\n}", { label: "try", detail: "try/catch block", type: "snippet" })
    ]
  };
}

// Custom linter for basic syntax errors
const basicLinter = linter((view) => {
  const diagnostics: readonly Diagnostic[] = [];
  const content = view.state.doc.toString();
  
  // Simple checks for syntax issues
  if (language === 'node') {
    // Check for common JavaScript syntax issues
    const unmatchedParentheses = (content.match(/\(/g) || []).length !== (content.match(/\)/g) || []).length;
    const unmatchedBraces = (content.match(/{/g) || []).length !== (content.match(/}/g) || []).length;
    const unmatchedBrackets = (content.match(/\[/g) || []).length !== (content.match(/\]/g) || []).length;
    
    if (unmatchedParentheses) {
      return [{
        from: 0,
        to: content.length,
        severity: "warning",
        message: "Unmatched parentheses detected"
      }];
    }
    if (unmatchedBraces) {
      return [{
        from: 0,
        to: content.length,
        severity: "warning",
        message: "Unmatched braces detected"
      }];
    }
    if (unmatchedBrackets) {
      return [{
        from: 0,
        to: content.length,
        severity: "warning",
        message: "Unmatched brackets detected"
      }];
    }
  }
  
  return diagnostics;
});

// Completion extensions
const completionExtensions: Record<string, Extension> = {
  node: autocompletion({ override: [jsCompletions] }),
  php: autocompletion({ override: [phpCompletions] }),
  go: autocompletion({ override: [goCompletions] })
};

// Base extensions for the editor
const baseExtensions: Extension[] = [
  lineNumbers(),
  highlightActiveLineGutter(),
  foldGutter(),
  highlightSpecialChars(),
  history(),
  drawSelection(),
  dropCursor(),
  EditorState.allowMultipleSelections.of(true),
  indentOnInput(),
  syntaxHighlighting(defaultHighlightStyle),
  bracketMatching(),
  closeBrackets(),
  autocompletion(),
  rectangularSelection(),
  crosshairCursor(),
  highlightActiveLine(),
  highlightSelectionMatches(),
  search({ top: true }), // Add search bar at the top
  basicLinter,
  ViewPlugin.fromClass(class {
    update(update: ViewUpdate) {
      if (update.docChanged) {
        onChange(update.state.doc.toString());
      }
    }
  }),
  keymap.of([
    indentWithTab,
    ...closeBracketsKeymap,
    ...defaultKeymap,
    ...searchKeymap,
    ...historyKeymap,
    ...foldKeymap,
    ...completionKeymap,
    ...lintKeymap
  ]),
  ...themeExtension
];

// Initialize editor on mount
onMount(() => {
  const langSupport = languageConfigs[language] || languageConfigs.node;
  const completionExt = completionExtensions[language] || completionExtensions.node;
  
  const initialState = EditorState.create({
    doc: code,
    extensions: [
      ...baseExtensions,
      langSupport,
      completionExt
    ]
  });

  view = new EditorView({
    state: initialState,
    parent: editorContainer,
    dispatch: (tr) => {
      view?.update([tr]);
      if (tr.docChanged && view) {
        onChange(view.state.doc.toString());
      }
    }
  });
});

// Track previous language, theme and code to detect changes
let prevLanguage = $state(language);
let prevTheme = $state(theme);
let prevCode = $state(code);

// Update editor when language or theme changes (but not when code changes)
$effect(() => {
  if (view && (prevLanguage !== language || prevTheme !== theme)) {
    const langSupport = languageConfigs[language] || languageConfigs.node;
    const completionExt = completionExtensions[language] || completionExtensions.node;
    const themeExt = theme === 'dark' ? [oneDark] : [];
    
    // Only update extensions, keep the same document content
    const newState = EditorState.create({
      doc: view.state.doc,  // Keep current document content
      extensions: [
        ...baseExtensions.slice(0, -1), // Remove the last extension (theme)
        ...themeExt,
        langSupport,
        completionExt
      ]
    });
    
    view.setState(newState);
    
    // Update previous values
    prevLanguage = language;
    prevTheme = theme;
  }
});

// Update document content when code prop changes (external update, e.g., language switch)
$effect(() => {
  if (view && prevCode !== code && code !== view.state.doc.toString()) {
    // Update document content
    view.dispatch({
      changes: { from: 0, to: view.state.doc.length, insert: code }
    });
    prevCode = code;
  } else if (view) {
    // Update prevCode to current doc content if they're the same
    prevCode = view.state.doc.toString();
  }
});

// Clean up
onDestroy(() => {
  view?.destroy();
  view = null;
});

// Expose methods to parent components
export function getEditorContent(): string {
  return view ? view.state.doc.toString() : code;
}

export function setEditorContent(content: string): void {
  code = content;
  if (view) {
    view.dispatch({
      changes: { from: 0, to: view.state.doc.length, insert: content }
    });
  }
}
</script>

<div bind:this={editorContainer} class="w-full h-full"></div>