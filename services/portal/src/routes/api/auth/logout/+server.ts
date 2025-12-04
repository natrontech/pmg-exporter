import { json } from '@sveltejs/kit';
import { logger } from '$lib/logger';

export async function POST({ cookies }) {
	logger.info('Clearing session cookies');
	try {
		const options = { path: '/' };
		cookies.delete('id_token', options);
		cookies.delete('access_token', options);
		return json({ success: true }, { status: 200 });
	} catch (err) {
		const errorMessage =
			err && typeof err === 'object' && 'message' in err ? err.message : String(err);
		logger.error({ error: errorMessage }, 'Error clearing session cookies');
		return json({ error: 'Failed to clear session' }, { status: 500 });
	}
}
