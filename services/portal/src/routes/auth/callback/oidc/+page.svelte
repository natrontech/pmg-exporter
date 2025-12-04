<script lang="ts">
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { userManager } from '$lib/client/oidc';
	import { logger } from '$lib/logger'; // Assuming logger can run client-side or is tree-shaken

	let status = 'Processing callback...';
	let error: string | null = null;

	onMount(async () => {
		if (!userManager) {
			error = 'OIDC User Manager is not available. Check browser context and configuration.';
			status = 'Error';
			logger.error(error);
			return;
		}

		try {
			logger.debug('Attempting signinCallback on client');
			const user = await userManager.signinCallback();

			if (user && user.access_token && user.id_token) {
				logger.info('SigninCallback successful, got user, access token and id token');
				status = 'Setting session...';

				const response = await fetch('/api/auth/session', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json'
					},
					body: JSON.stringify({
						id_token: user.id_token,
						access_token: user.access_token
					})
				});

				if (!response.ok) {
					const errorBody = await response.json().catch(() => ({}));
					const serverError = errorBody.error || 'Failed to set session cookie.';
					throw new Error(`Server error setting session: ${response.status} - ${serverError}`);
				}

				status = 'Redirecting...';
				await goto('/overview?toast=login_success', { replaceState: true });
			} else {
				throw new Error(
					'SigninCallback completed but did not return valid tokens (id_token and access_token required).'
				);
			}
		} catch (err: unknown) {
			const errorMessage = err instanceof Error ? err.message : String(err);
			logger.error({ error: errorMessage }, 'Error processing OIDC callback on client');
			error = errorMessage || 'An unknown error occurred during login.';
			status = 'Login Failed';
			// Optional: Redirect to home with error after a delay?
			// setTimeout(() => goto('/?error=callback_failed'), 5000);
		}
	});
</script>

<div class="flex min-h-screen flex-col items-center justify-center p-4">
	<h1 class="mb-4 text-2xl font-semibold">Authenticating...</h1>
	<p class="text-base-content/80 mb-2">{status}</p>

	{#if status !== 'Processing callback...' && status !== 'Setting session...' && status !== 'Redirecting...'}
		<!-- Show spinner only during active processing -->
	{:else}
		<span class="loading loading-dots loading-lg"></span>
	{/if}

	{#if error}
		<div class="bg-error text-error-content mt-4 max-w-md rounded-lg p-4">
			<p class="font-bold">Error:</p>
			<p class="text-sm">{error}</p>
			<p class="mt-2 text-xs">
				You will be redirected shortly, or you can <a href="/" class="link">return home now</a>.
			</p>
		</div>
	{/if}
</div>
