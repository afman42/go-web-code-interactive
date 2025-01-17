<script lang="ts">
import CodeMirror from "svelte-codemirror-editor";
import type { EditorView } from "@codemirror/view";
import { onDestroy} from "svelte";
import { javascript} from "@codemirror/lang-javascript"
import { langState } from "./utils/lang-state"
let view: EditorView;
let stdout = $state("Nothing");
let stderr = $state("Nothing");
let disabled = $state(false);
onDestroy(() => {
  view.destroy();
})
function onChange(e: CustomEvent){
  $langState.sampleDataLang[$langState.type][$langState.value] = e.detail
}
async function send(){
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
    stdout = res.errout != "" ? "Nothing" : $langState.type == "stq" ? JSON.parse(res.out.trim()) : res.out
    stderr = res.errout == "" ? "Nothing" : res.errout
  }

}
function onChangeRadio(event: Event){
  $langState.value = (event.target as HTMLInputElement).value
}
function onChangeType(event: Event){
  $langState.type = (event.target as HTMLInputElement).value
}
</script>

<main class="columns">
  <div class="column" style="margin:1rem;">
    <div style="margin-bottom:0.5rem;border-color:black; border-width: 2px;border-style: solid;">
      <CodeMirror bind:value={$langState.sampleDataLang[$langState.type][$langState.value]} readonly={disabled} on:change={(e) => onChange(e)} on:ready={(e) => view = e.detail} lang={javascript()}/> 
    </div> 
    <div>
        <button class="button is-light" disabled={disabled} onclick={send} type="button">Send</button> 
        <label class="radio">
        <input type="radio" value="node" onchange={onChangeRadio} checked={$langState.value == "node"}/>
          Node
        </label>
        <label class="radio">
          <input type="radio" value="php" onchange={onChangeRadio} checked={$langState.value == "php"} />
          PHP
        </label>
        <label class="radio">
          <input type="radio" value="go" onchange={onChangeRadio} checked={$langState.value == "go"} />
          Go
        </label>
        <label class="radio"> || </label>
        <label class="radio">
          <input type="radio" value="repl" onchange={onChangeType} checked={$langState.type == "repl"} />
          REPL
        </label>
        <label class="radio">
          <input type="radio" value="stq" onchange={onChangeType} checked={$langState.type == "stq"} />
          Simple Test Question
        </label>
      </div>
    </div>

  <div class="column">
    <div class="column content">

    {#if $langState.type == "repl"}
      StdOut <br />
      <blockquote>{stdout}</blockquote>
    {/if}
    {#if $langState.type == "stq"}
      Simple Test Question : change integer to string <br />
      Result <br />
      <blockquote>
        <input type="checkbox" value="stq1" checked={!stdout} /> Check change after int to string <br />
      </blockquote>
    {/if}
    </div>
    <div class="column content">
      StdErr <br />
      <blockquote>{stderr}</blockquote>
    </div>
  </div>
</main>

