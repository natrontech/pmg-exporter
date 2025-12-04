<script lang="ts">
	import { PlusCircle, Cable } from 'lucide-svelte';
	import type { LayoutData } from '../../routes/(app)/$types';

	// We might need data later if we display user-specific info here
	let { data } = $props<{ data: LayoutData }>();

	$effect(() => {
		console.log(data);
	});

	// Placeholder for active sessions data
	const activeSessions = $state([
		{ id: 'sess-123', name: 'My Main Project', url: '#' },
		{ id: 'sess-456', name: 'Quick Bugfix', url: '#' }
	]);

	function connectToSession(url: string) {
		// Logic to open the session URL, potentially in a new tab
		window.open(url, '_blank');
	}
</script>

<header class="border-base-300 bg-base-100 flex h-16 items-center justify-between border-b">
	<!-- Logo Container: Fixed width like sidebar, centered content, border -->
	<div
		class="border-base-300 flex h-full w-64 flex-shrink-0 items-center justify-center border-r px-4"
	>
		<a href="/overview" class="flex items-center gap-2">
			<img src="/logo/koda_logo.svg" alt="Koda Icon" class="h-8" />
		</a>
	</div>

	<!-- Right-side actions - Adjust padding -->
	<div class="flex flex-1 items-center justify-end gap-4 px-4 md:px-6">
		<div class="dropdown dropdown-end">
			<button tabindex="0" class="btn btn-ghost btn-circle" aria-label="Active Sessions">
				<Cable class="h-5 w-5" />
			</button>
			<ul
				class="dropdown-content menu menu-sm bg-base-100 rounded-box border-base-300 z-[1] mt-3 w-52 border p-2 shadow"
			>
				{#if activeSessions.length > 0}
					<li class="menu-title text-xs"><span>Active Sessions</span></li>
					{#each activeSessions as session (session.id)}
						<li>
							<button onclick={() => connectToSession(session.url)} class="text-xs">
								{session.name}
							</button>
						</li>
					{/each}
				{:else}
					<li><span class="text-xs italic">(No active sessions)</span></li>
				{/if}
			</ul>
		</div>
		<button class="btn btn-ghost btn-circle" aria-label="Create New Dev Environment">
			<PlusCircle class="h-5 w-5" />
		</button>
	</div>
</header>
