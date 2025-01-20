<script lang="ts">
//Component
import Toaster from "./component/Toaster.svelte"
//Lib
import CodeMirror from "svelte-codemirror-editor";
import type { EditorView } from "@codemirror/view";
import { onDestroy} from "svelte";
import { javascript } from "@codemirror/lang-javascript"
import { go } from "@codemirror/lang-go"
import { php, phpLanguage } from "@codemirror/lang-php"
import { langState } from "./utils/lang-state"
import { toasts } from './utils/toast';
let view: EditorView;
let stdout = $state("Nothing");
let stderr = $state("Nothing");
let disabled = $state(false);
let count = $state(0);
onDestroy(() => {
  view.destroy();
})
function onChange(e: CustomEvent){
  $langState.sampleDataLang[$langState.type][$langState.value] = e.detail
}
async function send(){
  toasts.info("Waiting Response",1000)
  disabled = true
  const payload = {
    "txt": $langState.sampleDataLang[$langState.type][$langState.value] as string,
    "lang":$langState.value,
    "type": $langState.type
  } as { [key: string]: string }
  let fetch = import("./utils/fetch"); 
  const res = await (await fetch).fetchApiPost<FetchData>(payload,"/") 
  if(res.statusCode == 200){
    disabled = false
    stderr = res.errout == "" ? "Nothing" : res.errout
    stdout = res.errout != "" ? "Nothing" : $langState.type == "stq" ? JSON.parse(res.out.trim())  : res.out
    if(stderr != "Nothing") {
      toasts.warning("Something Went Wrong",1000)
    }
    if(stdout != "Nothing" ){
      toasts.success("Success Response",1000)
    }
  }
}
function onChangeRadio(event: Event){
  $langState.value = (event.target as HTMLInputElement).value
  count++
  stdout = "Nothing"
  stderr = "Nothing" 
}
function onChangeType(event: Event){
  $langState.type = (event.target as HTMLInputElement).value
  count++
  stdout = "Nothing"
  stderr = "Nothing"
}
</script>

<div class="grid grid-cols-2 gap-4 sm:flex sm:flex-col min-sm:flex min-sm:flex-col">
  <div class="flex flex-col m-4">
    <div class="flex" style="width:100%;margin-bottom:0.5rem;border-color:black; border-width: 2px;border-style: solid;">
      {#key count}
        <CodeMirror class="w-full" bind:value={$langState.sampleDataLang[$langState.type][$langState.value]} readonly={disabled} on:change={(e) => onChange(e)} on:ready={(e) => view = e.detail} lang={javascript()} extensions={[go(),php({ baseLanguage: phpLanguage})]}/> 
      {/key}
    </div> 
    <div class="flex items-center md:flex md:items-center sm:flex sm:items-center min-sm:flex min-sm:flex-col">
        <button class="bg-red-500 flex py-2.5 px-3 min-sm:h-8 min-sm:hover:bg-white min-sm:hover:border-2 min-sm:hover:text-black min-sm:hover:border-red-500 min-sm:w-full md:px-1 md:py-2 md:text-sm min-sm:text-xs min-sm:items-center min-sm:justify-center text-white rounded-lg mr-1" disabled={disabled} onclick={send} type="button">Send</button> 
        <div class="flex gap-1 min-sm:justify-between">
          <div class="flex gap-2 md:flex md:items-center sm:flex sm:items-center sm:gap-1 min-sm:flex min-sm:items-center min-sm:gap-1">
          <label for="node" class="min-sm:text-sm min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="node" onchange={onChangeRadio} checked={$langState.value == "node"}/>Node
            </label>
          <label for="php" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="php" onchange={onChangeRadio} checked={$langState.value == "php"} />PHP
          </label>
          <label for="go" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
          <input type="radio" value="go" onchange={onChangeRadio} checked={$langState.value == "go"} />Go</label>
          <b> || </b>
          <label for="repl" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="repl" onchange={onChangeType} checked={$langState.type == "repl"} />REPL
          </label>
          <label for="stq" class="min-sm:text-xs min-sm:flex min-sm:gap-1 md:text-sm">
            <input type="radio" value="stq" onchange={onChangeType} checked={$langState.type == "stq"} />Simple Test Question
          </div>
        </div>
      </div>
    </div>

  <div class="flex flex-col ml-4">
    <div class="flex mt-3 flex-col">

    {#if $langState.type == "repl"}
      <h6 class="min-sm:text-sm md:text-sm">StdOut</h6>
      <blockquote class="min-sm:text-sm border-l-4 border-gray-500 my-2 py-4 pl-4 md:text-sm">{stdout}</blockquote>
    {/if}
    {#if $langState.type == "stq"}
      <h6 class="min-sm:text-sm md:text-sm">Simple Test Question : change integer to string</h6>
      <h6 class="min-sm:text-sm md:text-sm">Result</h6>
      <blockquote class="flex gap-1 min-sm:text-sm md:text-sm flex-start border-l-4 border-gray-500 my-2 py-4 pl-4">
        <input type="checkbox" value="stq1" checked={!stdout} disabled /> Check change after int to string
      </blockquote>
    {/if}
      <h6 class="min-sm:text-sm md:text-sm">StdErr</h6>
      <blockquote  class="border-l-4 md:text-sm min-sm:text-sm border-gray-500 my-2 py-4 pl-4">{stderr}</blockquote>
    </div>
  </div>
</div>

<Toaster />
