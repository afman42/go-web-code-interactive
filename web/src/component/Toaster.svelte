<script lang="ts">
	import { flip } from 'svelte/animate';
	import { fly } from 'svelte/transition';
	import { toasts, type IToast } from '../utils/toast';
  let themes = {
		error: '#E26D69',
		success: '#84C991',
		warning: '#f0ad4e',
		info: '#5bc0de'
	};
  const { stateToast } = toasts;
  //Should Writable to assign type guard
  const t = $derived($stateToast) as IToast[]
</script>

<div class="fixed top-0 right-0 mx-auto p-0 z-50 left-4 flex flex-col justify-start items-center pointer-events-none">
	{#each t as toast (toast.id)}
		<div
			animate:flip
			class="mb-2"
			style="background: {themes[toast.type]};flex: '0 0 auto';"
			transition:fly={{ y: 30 }}
		>
			<div class="p-3 block text-white font-medium text-xl min-sm:text-xs min-sm:font-normal sm:text-sm sm:font-normal">{toast.message}</div>
		</div>
	{/each}
</div>

