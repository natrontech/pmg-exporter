import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	// If the user is already logged in when visiting the root page,
	// redirect them to their profile or a dashboard.
	if (locals.user) {
		throw redirect(307, '/overview');
	}

	// No specific data needed for the login page itself,
	// the layout already provides providerName if needed.
	return {};
};
