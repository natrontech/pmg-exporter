import { UserManager, WebStorageStateStore } from 'oidc-client-ts';
import { env } from '$env/dynamic/public';
import { browser } from '$app/environment';

// Extract base authority from discovery URL if needed (usually not for v2.0 endpoint)
function getAuthority(url: string): string {
	// For v2.0 endpoints, the authority is usually the URL up to /v2.0
	const v2Marker = '/v2.0';
	const wellKnownMarker = '/.well-known/openid-configuration';
	if (url.includes(v2Marker)) {
		return url.substring(0, url.indexOf(v2Marker) + v2Marker.length);
	} else if (url.includes(wellKnownMarker)) {
		// Fallback for other issuer formats
		return url.substring(0, url.indexOf(wellKnownMarker));
	}
	return url; // Assume the URL is the authority if markers aren't present
}

const settings = {
	authority: getAuthority(env.PUBLIC_OIDC_ISSUER),
	metadataUrl: env.PUBLIC_OIDC_ISSUER.endsWith('/.well-known/openid-configuration')
		? env.PUBLIC_OIDC_ISSUER // If it already includes well-known path
		: `${getAuthority(env.PUBLIC_OIDC_ISSUER)}/.well-known/openid-configuration`, // Construct metadata URL
	client_id: env.PUBLIC_OIDC_CLIENT_ID,
	redirect_uri: `${env.PUBLIC_BASE_URL}/auth/callback/oidc`,
	post_logout_redirect_uri: env.PUBLIC_BASE_URL,
	response_type: 'code', // Use Authorization Code Flow with PKCE
	scope: 'openid profile email groups', // Adjust scopes as needed

	// Use browser session storage for state
	userStore: browser ? new WebStorageStateStore({ store: window.sessionStorage }) : undefined,
	stateStore: browser ? new WebStorageStateStore({ store: window.sessionStorage }) : undefined

	// Optional: Automatic silent renew (consider implications)
	// automaticSilentRenew: true,
	// silent_redirect_uri: `${PUBLIC_BASE_URL}/auth/silent-renew`, // Requires a silent renew page

	// Optional: Monitor session - checks iframe for login status changes at IdP
	// monitorSession: true,
	// checkSessionInterval: 10000, // milliseconds
};

// Create manager only in browser context
export const userManager = browser ? new UserManager(settings) : null;

// Function to trigger login
export async function login() {
	if (userManager) {
		await userManager.signinRedirect();
	} else {
		console.error('UserManager not available.');
		// Handle error appropriately, maybe redirect to an error page
	}
}

// Function to trigger logout
export async function logout() {
	if (userManager) {
		try {
			// 1. Clear the httpOnly cookie via server endpoint
			await fetch('/api/auth/logout', { method: 'POST' });
		} catch (error) {
			console.error('Failed to clear server session cookie:', error);
			// Decide if you want to proceed with IdP logout anyway
		}

		// 2. Trigger IdP signout and clear local OIDC state
		await userManager.signoutRedirect();
	} else {
		console.error('UserManager not available.');
	}
}
