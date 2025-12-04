import { tick } from 'svelte';
import { toast } from 'svelte-sonner';
import { replaceState } from '$app/navigation';
import { browser } from '$app/environment';

/**
 * Checks the current URL for a 'toast' search parameter and displays a corresponding toast message.
 * Removes the 'toast' parameter from the URL after displaying the message.
 * Should be called within onMount or equivalent client-side logic.
 */
export async function handleUrlToast(): Promise<void> {
	if (!browser) {
		// Ensure this only runs on the client-side
		return;
	}

	const params = new URLSearchParams(window.location.search);
	const toastParam = params.get('toast');

	if (toastParam) {
		// Map parameter values to toast messages and types
		const toastMap: Record<
			string,
			{ type: 'success' | 'error' | 'info' | 'warning'; message: string }
		> = {
			login_success: { type: 'success', message: 'Successfully logged in!' },
			logout_success: { type: 'success', message: 'Successfully logged out!' }
			// Add more mappings here as needed (e.g., different types)
			// 'profile_updated': { type: 'info', message: 'Profile updated.' },
			// 'action_failed': { type: 'error', message: 'Action failed. Please try again.' },
		};

		const toastInfo = toastMap[toastParam];

		if (toastInfo) {
			switch (toastInfo.type) {
				case 'success':
					toast.success(toastInfo.message);
					break;
				case 'error':
					toast.error(toastInfo.message);
					break;
				case 'info':
					toast.info(toastInfo.message);
					break;
				case 'warning':
					toast.warning(toastInfo.message);
					break;
			}
		} else {
			console.warn(`Unknown toast parameter value: ${toastParam}`);
		}

		// Remove the parameter from the URL to prevent re-triggering
		const url = new URL(window.location.href);
		url.searchParams.delete('toast');
		await tick(); // Ensure Svelte updates are processed if needed
		replaceState(url, {}); // Update URL state without a full page reload
	}
}
