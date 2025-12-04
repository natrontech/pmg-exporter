import { redirect } from '@sveltejs/kit';
import { logger } from '$lib/logger';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals }) => {
	logger.debug({ layout_locals_user: locals.user }, 'Checking locals.user in (app) layout load');
	if (!locals.user) {
		logger.info('User not authenticated, redirecting to login page (/)');
		throw redirect(307, '/');
	}

	return {
		user: locals.user
	};
};
