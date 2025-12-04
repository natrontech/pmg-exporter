<script lang="ts">
	import { onMount } from 'svelte';
	import { handleUrlToast } from '$lib/utils/toast';

	onMount(async () => {
		await handleUrlToast();
	});
</script>

<div>
	<h1>Overview</h1>
</div>
