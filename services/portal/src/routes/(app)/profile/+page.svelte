<script lang="ts">
	import type { PageData } from './$types';

	let { data } = $props<{ data: PageData }>();
</script>

<div class="prose">
	<h1>User Profile</h1>
	{#if data.profile}
		<div class="not-prose grid gap-4">
			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title">Personal Information</h2>
					<div class="space-y-2">
						<p><strong>Email:</strong> {data.profile.email}</p>
						<p><strong>Name:</strong> {data.profile.name || 'Not provided'}</p>
						<p><strong>Subject ID:</strong> {data.profile.sub}</p>
					</div>
				</div>
			</div>

			<div class="card bg-base-100 shadow">
				<div class="card-body">
					<h2 class="card-title">Authorization</h2>
					<div class="space-y-2">
						<div>
							<strong>Roles:</strong>
							{#if data.profile.roles && data.profile.roles.length > 0}
								<div class="mt-1 flex flex-wrap gap-1">
									{#each data.profile.roles as role (role)}
										<span class="badge badge-primary">{role}</span>
									{/each}
								</div>
							{:else}
								<span class="text-base-content/70">No roles assigned</span>
							{/if}
						</div>

						<div>
							<strong>Groups:</strong>
							{#if data.profile.groups && data.profile.groups.length > 0}
								<div class="mt-1 flex flex-wrap gap-1">
									{#each data.profile.groups as group (group)}
										<span class="badge badge-secondary">{group}</span>
									{/each}
								</div>
							{:else}
								<span class="text-base-content/70">No groups assigned</span>
							{/if}
						</div>
					</div>
				</div>
			</div>
		</div>
	{:else}
		<div class="flex items-center justify-center py-8">
			<span class="loading loading-spinner loading-lg"></span>
			<span class="ml-2">Loading user profile...</span>
		</div>
	{/if}
</div>
