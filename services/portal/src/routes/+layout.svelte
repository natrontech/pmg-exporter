<script lang="ts">
	import '../app.css';
	import type { Snippet } from 'svelte';
	import { Toaster } from 'svelte-sonner';

	let { children } = $props<{ children: Snippet /* data: LayoutData */ }>();
</script>

<svelte:head>
	<title>koda</title>
</svelte:head>

<main>
	{@render children()}
	<Toaster richColors position="bottom-center" />
</main>
