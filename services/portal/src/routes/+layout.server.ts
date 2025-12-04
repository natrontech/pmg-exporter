import { env } from '$env/dynamic/public';

export type LayoutData = {
	loggedIn: boolean;
	user: App.Locals['user'] | null;
	providerName: string;
};

export function load({ locals }): LayoutData {
	// Determine provider name from environment or default
	const providerName = env.PUBLIC_OIDC_PROVIDER_NAME || 'OIDC Provider';

	return {
		loggedIn: !!locals.user,
		user: locals.user ?? null,
		providerName: providerName
	};
}
