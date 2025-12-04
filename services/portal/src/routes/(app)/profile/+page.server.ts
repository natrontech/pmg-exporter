import type { PageServerLoad } from './$types';

export const load: PageServerLoad = ({ locals }) => {
	// The layout already ensures the user exists
	return {
		profile: locals.user
	};
};
