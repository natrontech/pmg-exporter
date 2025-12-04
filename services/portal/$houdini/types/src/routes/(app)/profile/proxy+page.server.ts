// @ts-nocheck
import type { PageServerLoad } from './$types';

export const load = ({ locals }: Parameters<PageServerLoad>[0]) => {
	// The layout already ensures the user exists
	return {
		profile: locals.user
	};
};
