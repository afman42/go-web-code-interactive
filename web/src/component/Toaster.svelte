<script lang="ts">
	import { flip } from 'svelte/animate';
	import { fly } from 'svelte/transition';
	import { toasts } from '../utils/toast';

	let themes = {
		error: '#E26D69',
		success: '#84C991',
		warning: '#f0ad4e',
		info: '#5bc0de'
	};
	const { state } = toasts;
</script>

<div class="notifications">
	{#each $state as toast (toast.id)}
		<div
			animate:flip
			class="toast"
			style="background: {themes[toast.type]};"
			transition:fly={{ y: 30 }}
		>
			<div class="content">{toast.message}</div>
		</div>
	{/each}
</div>

<style>
	.notifications {
		position: fixed;
		top: 0;
		left: 10px;
		right: 0;
		margin: 0 auto;
		padding: 0;
		z-index: 9999;
		display: flex;
		flex-direction: column;
		justify-content: flex-start;
		align-items: center;
		pointer-events: none;
	}

	.toast {
		flex: 0 0 auto;
		margin-bottom: 10px;
	}

	.content {
		padding: 10px;
		display: block;
		color: white;
		font-weight: 500;
    font-size: 20px;
	}
</style>
