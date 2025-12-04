<script lang="ts">
	import type { Snippet } from 'svelte';
	import type { LayoutData } from './$types';
	import Sidenavbar from '$lib/components/Sidenavbar.svelte';
	import Topbar from '$lib/components/Topbar.svelte';

	// Get user info from the root layout data
	let { children, data } = $props<{ children: Snippet; data: LayoutData }>();
</script>

<div class="flex h-screen flex-col">
	<!-- Top Navigation Bar -->
	<Topbar {data} />

	<div class="flex flex-1 overflow-hidden">
		<!-- Side Navigation Bar -->
		<div class="bg-base-100 border-base-300 w-64 flex-shrink-0 border-r">
			<Sidenavbar {data} />
		</div>

		<!-- Page content -->
		<main class="flex-1 overflow-y-auto p-4 md:p-8">
			{@render children()}
		</main>
	</div>
</div>
