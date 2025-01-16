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
  $langState.sampleDataLang[$langState.value] = e.detail
}
async function send(){
  disabled = true
  const payload = {
    "txt": $langState.sampleDataLang[$langState.value] as string,
    "lang":$langState.value
  } as { [key: string]: string }
  let fetch = import("./utils/fetch"); 
  const res = await (await fetch).fetchApiPost<FetchData>(payload,"/") 
  if(res.statusCode == 200){
    disabled = false
    stdout = res.errout != "" ? "Nothing" : res.out
    stderr = res.errout == "" ? "Nothing" : res.errout
  }

}
</script>

<main class="columns">
  <div class="column" style="margin:1rem;">
    <div style="margin-bottom:0.5rem;border-color:black; border-width: 2px;border-style: solid;">
      <CodeMirror bind:value={$langState.sampleDataLang[$langState.value]} readonly={disabled} on:change={(e) => onChange(e)} on:ready={(e) => view = e.detail} lang={javascript()}/> 
    </div> 
    <div>
        <button class="button is-light" disabled={disabled} onclick={send} type="button">Send</button> 
        <div class="select">
          <select bind:value={$langState.value} onchange={(e) => {
            let v = (e.target as HTMLInputElement).value
            $langState.value = v
          }}>
            <option value={"node"}>Node</option>
            <option value={"php"}>PHP<option>
          </select>
        </div>
      </div>
    </div>

  <div class="column">
    <div class="column content">
      StdOut <br />
      <blockquote>{stdout}</blockquote>
    </div>
    <div class="column content">
      StdErr <br />
      <blockquote>{stderr}</blockquote>
    </div>
  </div>
</main>

