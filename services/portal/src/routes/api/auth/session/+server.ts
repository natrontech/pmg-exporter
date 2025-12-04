import { json } from '@sveltejs/kit';
import { logger } from '$lib/logger';
import * as jose from 'jose'; // Import jose

export async function POST({ request, cookies, url }) {
	logger.info('POST /api/auth/session endpoint hit');
	try {
		const { id_token, access_token } = await request.json();

		if (!id_token || typeof id_token !== 'string' || id_token.trim() === '') {
			logger.warn('Missing or invalid id_token in request body');
			return json({ error: 'Invalid id_token provided' }, { status: 400 });
		}

		if (!access_token || typeof access_token !== 'string' || access_token.trim() === '') {
			logger.warn('Missing or invalid access_token in request body');
			return json({ error: 'Invalid access_token provided' }, { status: 400 });
		}

		let tokenExpiry = 3600;
		try {
			const decodedPayload = jose.decodeJwt(access_token);
			if (decodedPayload && typeof decodedPayload.exp === 'number') {
				const expiresInSeconds = decodedPayload.exp - Math.floor(Date.now() / 1000);
				if (expiresInSeconds > 0) {
					tokenExpiry = expiresInSeconds;
				}
			}
		} catch (e) {
			logger.warn('Could not decode access_token with jose.decodeJwt, using default expiry.', {
				error: e
			});
		}

		logger.debug('Setting session cookies via API');
		const cookieOptions = {
			path: '/',
			httpOnly: true,
			secure: url.protocol === 'https:',
			maxAge: tokenExpiry,
			sameSite: 'lax' as const
		};

		cookies.set('id_token', id_token, cookieOptions);
		cookies.set('access_token', access_token, cookieOptions);

		return json({ success: true }, { status: 200 });
	} catch (err) {
		const errorMessage =
			err && typeof err === 'object' && 'message' in err ? err.message : String(err);
		logger.error({ error: errorMessage }, 'Error setting session cookie');
		return json({ error: 'Failed to set session' }, { status: 500 });
	}
}
