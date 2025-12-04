// src/hooks.server.ts
import { getAuthUser } from '$lib/auth';
import { setSession } from '$houdini';
import type { Handle } from '@sveltejs/kit';

export const handle: Handle = async ({ event, resolve }) => {
	const user = await getAuthUser(event.cookies);

	if (user) {
		event.locals.user = user;
		setSession(event, { user });
	}

	return resolve(event);
};
