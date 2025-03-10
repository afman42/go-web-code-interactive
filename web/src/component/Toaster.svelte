<script lang="ts">
  import { fade,slide } from 'svelte/transition';
	import { useToast, type IToastIDX } from '../utils/toast.svelte';
  let propToast: IToastIDX = $props()
  let themes = {
		error: '#E26D69',
		success: '#84C991',
		warning: '#f0ad4e',
		info: '#5bc0de'
	};
  const t = useToast()
  $effect(() => {
    if (t.state.length > 0) {
      const timer = setTimeout(() => {
        //Remove current idx
        t.state.splice(propToast.idx,1)
      }, t.state[propToast.idx].timeout);
      return () => {
        clearTimeout(timer);
      };
    }
  });
</script>

<div
  in:fade
  class="mb-2"
  style="background: {themes[propToast.type]};flex: '0 0 auto';"
  out:slide={{ axis: "y"}}
>
  <div class="p-3 block text-white font-medium text-xl min-sm:text-xs min-sm:font-normal sm:text-sm sm:font-normal">{propToast.message}</div>
</div>
