<script lang="ts">
	import type { PageData } from './$types';
	import { login } from '$lib/client/oidc'; // Import client-side login function

	// Interface to receive providerName from the layout
	interface ExtendedPageData extends PageData {
		providerName: string;
	}

	let { data } = $props<{ data: ExtendedPageData }>();

	// Function to get a provider identifier
	function getProviderIdentifier(name: string): string | null {
		const lowerName = name.toLowerCase();
		if (lowerName.includes('microsoft') || lowerName.includes('azure')) {
			return 'microsoft';
		}
		if (lowerName.includes('google')) {
			return 'google';
		}
		if (lowerName.includes('github')) {
			return 'github';
		}
		if (lowerName.includes('zitadel')) {
			return 'zitadel';
		}
		if (lowerName.includes('authentik')) {
			return 'authentik';
		}
		if (lowerName.includes('keycloak')) {
			return 'keycloak';
		}
		return null;
	}

	const providerIdentifier = getProviderIdentifier(data.providerName);
</script>

<div class="grid min-h-screen grid-cols-1 md:grid-cols-2">
	<!-- Left Side (Background/Graphic) -->
	<div
		class="from-primary to-secondary hidden items-center justify-center bg-gradient-to-br p-10 md:flex"
	>
		<div
			class="absolute inset-0 right-1/2 left-0 bg-[url('/img/login_background.png')] bg-cover bg-center"
		></div>
	</div>

	<!-- Right Side (Login Prompt Only) -->
	<div class="bg-base-100 flex flex-col items-center justify-center p-8">
		<img src="/logo/koda_logo.svg" alt="Koda Icon" class="mb-8 h-16" />

		<div class="w-full max-w-md text-center">
			<h1 class="mb-2 text-3xl font-bold">Sign In</h1>
			<p class="text-base-content/70 mb-6">Use your company account to continue.</p>
			<button type="button" onclick={login} class="btn btn-primary btn-block">
				{#if providerIdentifier}
					<img
						src={`/logo/providers/${providerIdentifier}.png`}
						alt={`${data.providerName} logo`}
						class="mr-2 h-5 w-5"
					/>
				{/if}
				Sign In with {data.providerName}
			</button>
			<p class="text-base-content/50 mt-4 text-xs">
				You will be redirected to {data.providerName} to authenticate.
			</p>
		</div>
	</div>
</div>
