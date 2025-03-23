<script lang="ts">
//Component
import ToasterContainer from "./component/ToasterContainer.svelte"
// AutoCompletion Just Javascript not another langugage, just add CompletionContext
// ref: https://grok.com/share/bGVnYWN5_9e2babff-5a19-44e7-8cdd-12b1461cf310
//Lib
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
  highlightActiveLineGutter
} from "@codemirror/view";
import {
  defaultHighlightStyle,
  syntaxHighlighting,
  indentOnInput,
  bracketMatching,
  foldGutter,
  foldKeymap,
  type LanguageSupport
} from "@codemirror/language";
import {
  defaultKeymap,
  history,
  historyKeymap,

  indentWithTab

} from "@codemirror/commands";
import {
  searchKeymap,
  highlightSelectionMatches
} from "@codemirror/search";
import {
  autocompletion,
  completionKeymap,
  closeBrackets,
  closeBracketsKeymap,
  type CompletionContext,
  type CompletionResult,
} from "@codemirror/autocomplete";
import { lintKeymap } from "@codemirror/lint";
import { javascript } from "@codemirror/lang-javascript";
import { go } from "@codemirror/lang-go";
import { php } from "@codemirror/lang-php";
// lib state
import { langState } from "./utils/lang-state.svelte"
import { useToast } from './utils/toast.svelte';
let stdout = $state("Nothing");
let stderr = $state("Nothing");
let disabled = $state(false);
const toast = useToast()
// Language configuration map
const languageConfigs = {
  node: javascript(),
  php: php(),
  go: go()
} as Record<string,LanguageSupport>;
// Reactive states
let view: EditorView | null = null;
let editorContainer: HTMLDivElement;
let currentLang = $state(languageConfigs.node) as any;
let editorValue = $state("");
let prevLang = $state(langState.value);
let prevType = $state(langState.type);

// Custom autocompletion for PHP
function phpCompletions(context: CompletionContext): CompletionResult | null {
  const word = context.matchBefore(/\w*/);
  if (!word || word.from === word.to && !context.explicit) return null;

  const phpKeywords = [
    "echo", "print", "if", "else", "foreach", "function", "return",
    "class", "public", "private", "protected", "namespace", "use"
  ];
  const phpFunctions = [
    "array", "strlen", "str_replace", "explode", "implode", "isset"
  ];

  return {
    from: word.from,
    options: [
      ...phpKeywords.map(kw => ({ label: kw, type: "keyword" })),
      ...phpFunctions.map(fn => ({ label: fn, type: "function" }))
    ]
  };
}

// Custom autocompletion for Go
function goCompletions(context: CompletionContext): CompletionResult | null {
  const word = context.matchBefore(/\w*/);
  if (!word || word.from === word.to && !context.explicit) return null;

  const goKeywords = [
    "func", "var", "const", "if", "else", "for", "range", "return",
    "struct", "interface", "package", "import", "type"
  ];
  const goBuiltins = [
    "println", "print", "len", "cap", "make", "new", "append"
  ];

  return {
    from: word.from,
    options: [
      ...goKeywords.map(kw => ({ label: kw, type: "keyword" })),
      ...goBuiltins.map(fn => ({ label: fn, type: "function" }))
    ]
  };
}

// Language-specific completion extensions
const completionExtensions = {
  node: autocompletion(), // Built-in for JavaScript
  php: autocompletion({ override: [phpCompletions] }),
  go: autocompletion({ override: [goCompletions] })
} as Record<string,Extension>;

// Base extensions (shared across all configurations)
const baseExtensions: Extension[] = [
  lineNumbers(),
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
  highlightActiveLineGutter(),
  highlightSelectionMatches(),
  keymap.of([
    indentWithTab,
    ...closeBracketsKeymap,
    ...defaultKeymap,
    ...searchKeymap,
    ...historyKeymap,
    ...foldKeymap,
    ...completionKeymap,
    ...lintKeymap
  ])
];
// Initialize editor on mount
onMount(() => {
  const initialState = EditorState.create({
    doc: langState.sampleDataLang[langState.type][langState.value] || "",
    extensions: [...baseExtensions, currentLang, completionExtensions[langState.value]]
  });

  view = new EditorView({
    state: initialState,
    parent: editorContainer,
    dispatch: (tr) => {
      view?.update([tr]);
      if (tr.docChanged) {
        editorValue = view?.state.doc.toString() || "";
        langState.sampleDataLang[langState.type][langState.value] = editorValue;
      }
    }
  });

  editorValue = view.state.doc.toString();
});

// Effect for language/type changes
$effect(() => {
  if (langState.value !== prevLang || langState.type !== prevType) {
    currentLang = languageConfigs[langState.value] || languageConfigs.node;
    editorValue = langState.sampleDataLang[langState.type][langState.value] || "";
    
    if (view) {
      const newState = EditorState.create({
        doc: editorValue,
        extensions: [...baseExtensions, currentLang, completionExtensions[langState.value]]
      });
      view.setState(newState);
    }
    
    prevLang = langState.value;
    prevType = langState.type;
  }
});
async function send(){
  toast.info("Waiting Response",3000)
  disabled = true
  const payload = {
    "txt": langState.sampleDataLang[langState.type][langState.value],
    "lang":langState.value,
    "type": langState.type
  } as { [key: string]: string }
  let fetch = import("./utils/fetch"); 
  try {
    const res = await (await fetch).fetchApiPost<FetchData>(payload,"/")
    if(res.statusCode == 200){
        disabled = false
        stderr = res.errout.trim().length > 0 ? res.errout : "Nothing"
        stdout = res.out.trim().length > 0 && stderr == "Nothing" ? res.out : "Nothing"
        if(langState.type == "stq") {
          stderr = res.errout.trim().length > 0 ? res.errout : "Nothing"
          stdout = res.out.trim().length > 0 ? JSON.parse(res.out.trim()) : "Nothing"
        }
        if(stderr != "Nothing") {
          toast.warning("Something Went Wrong",1000)
        }
        if(stdout != "Nothing" ){
          toast.success("Success Response",1000)
        }
      }
  } catch (error: unknown) {
    const parseMessage = JSON.parse(error as string)
    if(parseMessage.statusCode == 400) toast.error(parseMessage.message,1000)
    disabled = false
    stdout = "Nothing"
    stderr = "Nothing" 
  }
}
function onChangeRadio(event: Event){
  langState.value = (event.target as HTMLInputElement).value
  stdout = "Nothing"
  stderr = "Nothing" 
}
function onChangeType(event: Event){
  langState.type = (event.target as HTMLInputElement).value
  stdout = "Nothing"
  stderr = "Nothing"
}

// Clean up
onDestroy(() => {
  view?.destroy();
  view = null;
});

</script>

<div class="grid grid-cols-2 gap-4 sm:flex sm:flex-col min-sm:flex min-sm:flex-col">
  <div class="flex flex-col m-4">
    <div class="flex w-full border-black border-2 border-solid mb-2">
      <div bind:this={editorContainer} class="w-full"></div>
    </div> 
    <div class="flex items-center md:flex md:items-center sm:flex sm:items-center min-sm:flex min-sm:flex-col">
        <button class="bg-red-500 flex py-2.5 px-3 min-sm:h-8 min-sm:hover:bg-white min-sm:hover:border-2 min-sm:hover:text-black min-sm:hover:border-red-500 min-sm:w-full md:px-1 md:py-2 md:text-sm min-sm:text-xs min-sm:items-center min-sm:justify-center text-white rounded-lg mr-1" disabled={disabled} onclick={send} type="button">Send</button> 
        <div class="flex gap-1 min-sm:justify-between">
          <div class="flex gap-2 md:flex md:items-center sm:flex sm:items-center sm:gap-1 min-sm:flex min-sm:items-center min-sm:gap-1">
          <label for="node" class="min-sm:text-sm min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="node" onchange={onChangeRadio} checked={langState.value == "node"}/>Node
            </label>
          <label for="php" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="php" onchange={onChangeRadio} checked={langState.value == "php"} />PHP
          </label>
          <label for="go" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
          <input type="radio" value="go" onchange={onChangeRadio} checked={langState.value == "go"} />Go</label>
          <b> || </b>
          <label for="repl" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="repl" onchange={onChangeType} checked={langState.type == "repl"} />REPL
          </label>
          <label for="stq" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="stq" onchange={onChangeType} checked={langState.type == "stq"} />Simple Test Question
          </div>
        </div>
      </div>
    </div>

  <div class="flex flex-col ml-4">
    <div class="flex mt-3 flex-col">

    {#if langState.type == "repl"}
      <h6 class="min-sm:text-sm md:text-sm">StdOut</h6>
      <blockquote class="min-sm:text-sm border-l-4 border-gray-500 my-2 py-4 pl-4 md:text-sm">{stdout}</blockquote>
    {/if}
    {#if langState.type == "stq"}
      <h6 class="min-sm:text-sm md:text-sm">Simple Test Question : change integer to string</h6>
      <h6 class="min-sm:text-sm md:text-sm">Result</h6>
      <blockquote class="flex gap-1 min-sm:text-sm md:text-sm flex-start border-l-4 border-gray-500 my-2 py-4 pl-4">
        <input type="checkbox" value="stq1" checked={!stdout} disabled /> Check change after int to string
      </blockquote>
    {/if}
      <h6 class="min-sm:text-sm md:text-sm">StdErr</h6>
      <blockquote class="border-l-4 md:text-sm min-sm:text-sm border-gray-500 my-2 py-4 pl-4">{stderr}</blockquote>
    </div>
  </div>
</div>

<ToasterContainer />
