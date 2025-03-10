<script lang="ts">
	import { useToast } from '../utils/toast.svelte';
  import Toaster from "./Toaster.svelte"
  const t = useToast()
</script>

<div class="fixed top-0 right-0 mx-auto p-0 z-50 left-4 flex flex-col justify-start items-center pointer-events-none">
	{#each t.state as toast, i (toast.id)}
    <Toaster id={toast.id} idx={i} timeout={toast.timeout} message={toast.message} type={toast.type} />
	{/each}
</div>
