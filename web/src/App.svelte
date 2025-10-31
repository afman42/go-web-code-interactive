<script lang="ts">
// Component
import ToasterContainer from "./component/ToasterContainer.svelte"
import ExecuteButton from "./component/ExecuteButton.svelte"
import Editor from "./component/Editor.svelte"
// Lib state
import { langState } from "./utils/lang-state.svelte"
import { useToast } from './utils/toast.svelte';
import { validateUserCode } from './utils/validation';
import { TOAST_DURATION_MEDIUM, TOAST_DURATION_SHORT, SUPPORTED_LANGUAGES, EXECUTION_TYPES, THEMES } from './constants';

// Using Svelte 5 runes for state management
let stdout = $state("Nothing");
let stderr = $state("Nothing");
let isLoading = $state(false);
let editorValue = $state("");
let editorInitialized = $state(false);
let currentLangValue = $state(langState.value);
let currentTypeValue = $state(langState.type);
let editorRef: Editor | null = null;
let theme: 'light' | 'dark' = $state(THEMES.LIGHT);

// Reactive values using $derived
const isExecuting = $derived(isLoading);
const canExecute = $derived(!isLoading);

// Initialize editor value on first run
$effect(() => {
  if (!editorInitialized) {
    editorValue = langState.sampleDataLang[langState.type][langState.value];
    editorInitialized = true;
  }
});

// Update all values when language or type changes
$effect(() => {
  currentLangValue = langState.value;
  currentTypeValue = langState.type;
  const newCode = langState.sampleDataLang[langState.type][langState.value] || "";
  if (editorValue !== newCode) {  // Only update if different
    editorValue = newCode;
    if (editorRef) {
      editorRef.setEditorContent(newCode);
    }
  }
});

const toast = useToast();

async function send() {
  // Get current content from editor
  if (editorRef) {
    editorValue = editorRef.getEditorContent();
  }

  // Client-side validation
  const validation = validateUserCode(editorValue);
  if (!validation.isValid) {
    toast.error(`Validation failed: ${validation.errors.join(', ')}`, TOAST_DURATION_MEDIUM);
    return;
  }
  
  toast.info("Executing code...", TOAST_DURATION_SHORT);
  isLoading = true;
  
  // Use direct values from langState to ensure we have current values
  const payload = {
    "txt": editorValue,
    "lang": langState.value,
    "type": langState.type
  };
  
  try {
    const fetchModule = await import("./utils/fetch"); 
    const res = await fetchModule.fetchApiPost<FetchData>(payload, "/");
    
    if (res.statusCode === 200) {
      isLoading = false;
      stderr = res.errout.trim().length > 0 ? res.errout : "Nothing";
      stdout = res.out.trim().length > 0 && stderr === "Nothing" ? res.out : "Nothing";
      
      if (langState.type === EXECUTION_TYPES.STQ) {
        stderr = res.errout.trim().length > 0 ? res.errout : "Nothing";
        stdout = res.out.trim().length > 0 ? JSON.parse(res.out.trim()) : "Nothing";
      }
      
      if (stderr !== "Nothing") {
        toast.warning("Code executed with errors", TOAST_DURATION_SHORT);
      }
      if (stdout !== "Nothing") {
        toast.success("Code executed successfully", TOAST_DURATION_SHORT);
      }
    }
  } catch (error: unknown) {
    let errorMessage = "An error occurred";
    if (error instanceof Error) {
      errorMessage = error.message;
    } else if (typeof error === 'string') {
      try {
        const parsedError = JSON.parse(error);
        errorMessage = parsedError.message || error;
      } catch {
        errorMessage = error;
      }
    }
    
    toast.error(`Execution failed: ${errorMessage}`, TOAST_DURATION_MEDIUM);
    isLoading = false;
    stdout = "Nothing";
    stderr = "Nothing"; 
  }
}

function onChangeRadio(event: Event) {
  const target = event.target as HTMLInputElement;
  langState.value = target.value;
  stdout = "Nothing";
  stderr = "Nothing"; 
}

function onChangeType(event: Event) {
  const target = event.target as HTMLInputElement;
  langState.type = target.value;
  stdout = "Nothing";
  stderr = "Nothing";
}

function toggleTheme() {
  theme = theme === THEMES.LIGHT ? THEMES.DARK : THEMES.LIGHT;
}

function onEditorChange(value: string) {
  editorValue = value;
  // Update the langState with the current editor value using current lang and type
  langState.sampleDataLang[langState.type][langState.value] = value;
}
</script>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-6 p-4 max-w-7xl mx-auto">
  <!-- Editor Section -->
  <div class="flex flex-col">
    <div class="border-2 border-gray-300 dark:border-gray-600 rounded-lg overflow-hidden shadow-lg">
      <div class="h-[500px]">
        <Editor 
          bind:this={editorRef}
          language={currentLangValue}
          code={langState.sampleDataLang[currentTypeValue][currentLangValue]}
          onChange={onEditorChange}
          theme={theme}
        />
      </div>
    </div>
    
    <!-- Controls -->
    <div class="flex flex-wrap items-center gap-3 mt-4 p-3 bg-gray-100 dark:bg-gray-800 rounded-lg">
      <ExecuteButton 
        isLoading={isLoading}
        canExecute={canExecute}
        onClick={send}
      />
      
      <button 
        class="bg-gray-700 hover:bg-gray-800 text-white flex items-center justify-center py-2.5 px-4 rounded-lg transition-colors"
        onclick={toggleTheme}
        type="button"
        aria-label="Toggle theme"
      >
        {theme === THEMES.LIGHT ? '🌙 Dark Mode' : '☀️ Light Mode'}
      </button>
      
      <div class="flex flex-wrap gap-4 ml-auto">
        <!-- Language Selection -->
        <div class="flex items-center gap-2">
          <label class="flex items-center gap-1 cursor-pointer">
            <input 
              type="radio" 
              value={SUPPORTED_LANGUAGES.NODE} 
              onchange={onChangeRadio} 
              checked={langState.value === SUPPORTED_LANGUAGES.NODE}
            />
            <span  class="dark:text-white">Node</span>
          </label>
          <label class="flex items-center gap-1 cursor-pointer">
            <input 
              type="radio" 
              value={SUPPORTED_LANGUAGES.PHP} 
              onchange={onChangeRadio} 
              checked={langState.value === SUPPORTED_LANGUAGES.PHP}
            />
            <span  class="dark:text-white">PHP</span>
          </label>
          <label class="flex items-center gap-1 cursor-pointer">
            <input 
              type="radio" 
              value={SUPPORTED_LANGUAGES.GO} 
              onchange={onChangeRadio} 
              checked={langState.value === SUPPORTED_LANGUAGES.GO}
            />
            <span  class="dark:text-white">Go</span>
          </label>
        </div>
        
        <div class="text-gray-500 dark:text-gray-400 mx-2">|</div>
        
        <!-- Execution Type Selection -->
        <div class="flex items-center gap-2">
          <label class="flex items-center gap-1 cursor-pointer">
            <input 
              type="radio" 
              value={EXECUTION_TYPES.REPL} 
              onchange={onChangeType} 
              checked={langState.type === EXECUTION_TYPES.REPL}
            />
            <span  class="dark:text-white">REPL</span>
          </label>
          <label class="flex items-center gap-1 cursor-pointer">
            <input 
              type="radio" 
              value={EXECUTION_TYPES.STQ} 
              onchange={onChangeType} 
              checked={langState.type === EXECUTION_TYPES.STQ}
            />
            <span class="dark:text-white">STQ</span>
          </label>
        </div>
      </div>
    </div>
  </div>

  <!-- Output Section -->
  <div class="flex flex-col">
    <div class="bg-gray-50 dark:bg-gray-800 p-4 rounded-lg shadow-lg h-full">
      <h5 class="text-lg font-semibold mb-4 text-gray-800 dark:text-gray-200">Output</h5>
      
      {#if langState.type === EXECUTION_TYPES.STQ}
        <div class="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded border border-blue-200 dark:border-blue-800">
          <h6 class="font-medium text-blue-800 dark:text-blue-200 mb-1">Simple Test Question</h6>
          <p class="text-sm text-blue-600 dark:text-blue-300">Change integer to string</p>
        </div>
      {/if}
      
      <!-- StdOut -->
      <div class="mb-6">
        <div class="flex items-center justify-between mb-2">
          <h6 class="font-medium text-gray-700 dark:text-gray-300">StdOut</h6>
          {#if stdout !== "Nothing"}
            <span class="text-xs bg-green-100 text-green-800 px-2 py-1 rounded dark:bg-green-900/30 dark:text-green-300">
              Success
            </span>
          {/if}
        </div>
        <div class="min-h-[100px] max-h-40 overflow-auto p-3 bg-white dark:text-white dark:bg-gray-700/50 rounded border border-gray-200 dark:border-gray-600 text-sm font-mono whitespace-pre-wrap">
          {stdout === "Nothing" ? "No output" : stdout}
        </div>
      </div>
      
      <!-- StdErr -->
      <div>
        <div class="flex items-center justify-between mb-2">
          <h6 class="font-medium text-gray-700 dark:text-gray-300">StdErr</h6>
          {#if stderr !== "Nothing"}
            <span class="text-xs bg-red-100 text-red-800 px-2 py-1 rounded dark:bg-red-900/30 dark:text-red-300">
              Error
            </span>
          {/if}
        </div>
        <div class="min-h-[100px] max-h-40 overflow-auto p-3 bg-white dark:bg-gray-700/50 rounded border border-gray-200 dark:border-gray-600 text-sm font-mono whitespace-pre-wrap text-red-600 dark:text-red-400">
          {stderr === "Nothing" ? "No errors" : stderr}
        </div>
      </div>
    </div>
  </div>
</div>

<ToasterContainer />
