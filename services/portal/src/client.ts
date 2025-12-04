import { HoudiniClient } from '$houdini';
import type { Session } from '$houdini';
import { env } from '$env/dynamic/public';
import { browser } from '$app/environment';
import { toast } from 'svelte-sonner';

export default new HoudiniClient({
	url: env.PUBLIC_GRAPHQL_ENDPOINT,
	fetchParams({ session }: { session?: Session | null | undefined }) {
		const accessToken = session?.user?.access_token;

		const headers: Record<string, string> = {};

		if (accessToken && accessToken.trim() !== '') {
			headers.Authorization = `Bearer ${accessToken}`;
		}

		return {
			headers,
			credentials: 'include',
			signal: AbortSignal.timeout(10000)
		};
	},
	plugins: [
		() => ({
			catch: async (ctx, { error }) => {
				if (browser) {
					const errorMessage = error instanceof Error ? error.message : String(error);

					if (
						errorMessage.includes('fetch') ||
						errorMessage.includes('timeout') ||
						errorMessage.includes('ECONNREFUSED')
					) {
						toast.error('Backend service is temporarily unavailable');
					} else {
						toast.error('Something went wrong. Please try again.');
					}
				}

				console.error('GraphQL Error:', error);
			}
		})
	]
});
