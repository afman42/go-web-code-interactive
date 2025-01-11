<script lang="ts">
import CodeMirror from "svelte-codemirror-editor";
import type { EditorView } from "@codemirror/view";
import { onDestroy} from "svelte";
import { javascript} from "@codemirror/lang-javascript"
let view: EditorView;
let value = $state("console.log('hello world');\n\n\n\n\n\n\n\n");
let stdout = $state("Nothing");
let stderr = $state("Nothing");
let disabled = $state(false);
onDestroy(() => {
  view.destroy();
})
function onChange(e: CustomEvent){
  value = e.detail
}
async function send(){
  disabled = true
  const payload = {
    "txt": value
  }
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
      <CodeMirror bind:value on:change={(e) => onChange(e)} on:ready={(e) => view = e.detail} lang={javascript()}/> </div> <div> <button class="button is-light" disabled={disabled} onclick={send} type="button">Send</button> </div>
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

