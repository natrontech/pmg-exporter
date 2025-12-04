<script lang="ts">
	import type { PageData } from '$houdini/types/src/routes/(app)/namespaces/$houdini.js';

	let { data }: { data: PageData } = $props();
	let { Namespace } = $derived(data);
</script>

<div>
	<h1>Namespaces</h1>

	{#if $Namespace.fetching}
		<div class="flex justify-center py-8">
			<span class="loading loading-spinner loading-lg"></span>
		</div>
	{:else if $Namespace.data}
		{#each $Namespace.data.namespaces ?? [] as namespace (namespace.name)}
			<div class="card bg-base-100 mb-4 shadow-xl">
				<div class="card-body">
					<h2 class="card-title">{namespace.name}</h2>
				</div>
			</div>
		{/each}
	{:else}
		<div class="py-8 text-center">
			<p class="text-base-content/70">No namespaces found</p>
		</div>
	{/if}
</div>
