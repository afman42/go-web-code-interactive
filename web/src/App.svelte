<script lang="ts">
//Component
import Toaster from "./component/Toaster.svelte"
//Lib
import CodeMirror from "svelte-codemirror-editor";
import type { EditorView } from "@codemirror/view";
import { onDestroy} from "svelte";
import { javascript} from "@codemirror/lang-javascript"
import { go } from "@codemirror/lang-go"
import { php } from "@codemirror/lang-php"
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

<div class="grid grid-cols-2 gap-4 md:max-xl:flex-col md:max-xl:flex sm:max-xl:flex sm:max-xl:flex-col min-sm:max-xl:flex min-sm:max-xl:flex-col">
  <div class="flex flex-col" style="margin:1rem;">
    <div class="flex" style="width:100%;margin-bottom:0.5rem;border-color:black; border-width: 2px;border-style: solid;">
      {#key count}
        <CodeMirror class="w-full" bind:value={$langState.sampleDataLang[$langState.type][$langState.value]} readonly={disabled} on:change={(e) => onChange(e)} on:ready={(e) => view = e.detail} lang={javascript()} extensions={[go(),php()]}/> 
      {/key}
    </div> 
    <div class="flex items-center min-sm:max-xl:flex min-sm:max-xl:flex-col">
        <button class="bg-red-500 flex py-2.5 px-3 text-white rounded-lg mr-1" disabled={disabled} onclick={send} type="button">Send</button> 
        <div class="flex gap-1">
          <div class="flex gap-1 min-sm:max-xl:flex min-sm:max-xl:items-center min-sm:max-xl:gap-1">
          <label for="node" class="min-sm:max-xl:text-sm min-sm:max-xl:flex min-sm:max-xl:gap-1">
            <input type="radio" value="node" onchange={onChangeRadio} checked={$langState.value == "node"}/>Node
            </label>
          <label for="php" class="min-sm:max-xl:text-xs min-sm:max-xl:flex min-sm:max-xl:gap-1">
            <input type="radio" value="php" onchange={onChangeRadio} checked={$langState.value == "php"} />PHP
          </label>
          <label for="go" class="min-sm:max-xl:text-xs min-sm:max-xl:flex min-sm:max-xl:gap-1">
          <input type="radio" value="go" onchange={onChangeRadio} checked={$langState.value == "go"} />Go</label>
          <b> || </b>
          <label for="repl" class="min-sm:max-xl:text-xs min-sm:max-xl:flex min-sm:max-xl:gap-1">
            <input type="radio" value="repl" onchange={onChangeType} checked={$langState.type == "repl"} />REPL
          </label>
          <label for="stq" class="min-sm:max-xl:text-xs min-sm:max-xl:flex min-sm:max-xl:gap-1">
            <input type="radio" value="stq" onchange={onChangeType} checked={$langState.type == "stq"} />Simple Test Question
          </div>
        </div>
      </div>
    </div>

  <div class="flex flex-col md:max-xl:ml-4 sm:max-xl:ml-4 min-sm:max-xl:ml-4">
    <div class="flex mt-3 flex-col">

    {#if $langState.type == "repl"}
      <h6 class="min-sm:max-xl:text-sm">StdOut</h6>
      <blockquote class="min-sm:max-xl:text-sm border-l-4 border-gray-500 my-2 py-4 pl-4">{stdout}</blockquote>
    {/if}
    {#if $langState.type == "stq"}
      <h6 class="min-sm:max-xl:text-sm">Simple Test Question : change integer to string</h6>
      <h6 class="min-sm:max-xl:text-sm">Result</h6>
      <blockquote class="flex gap-1 min-sm:max-xl:text-sm flex-start border-l-4 border-gray-500 my-2 py-4 pl-4">
        <input type="checkbox" value="stq1" checked={!stdout} disabled /> Check change after int to string
      </blockquote>
    {/if}
      <h6 class="min-sm:max-xl:text-sm">StdErr</h6>
      <blockquote  class="border-l-4 min-sm:max-xl:text-sm border-gray-500 my-2 py-4 pl-4">{stderr}</blockquote>
    </div>
  </div>
</div>

<Toaster />
