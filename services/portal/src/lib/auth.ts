// src/lib/auth.ts
import { logger } from './logger';
import { env } from '$env/dynamic/public';
import * as jose from 'jose';
import { URL } from 'url'; // Node.js URL module

export interface UserClaims {
	sub: string;
	email: string;
	name?: string;
	roles: string[];
	groups: string[];
	id_token: string;
	access_token: string;
}

// Function to get the JWKS URI from the issuer URL
async function getJwksUri(issuerUrl: string): Promise<URL> {
	try {
		const issuer = new URL(issuerUrl);
		const response = await fetch(new URL('/.well-known/openid-configuration', issuer).toString());
		if (!response.ok) {
			throw new Error(`Failed to fetch OIDC discovery document: ${response.statusText}`);
		}
		const discoveryDoc = await response.json();
		if (!discoveryDoc.jwks_uri) {
			throw new Error('jwks_uri not found in OIDC discovery document');
		}
		return new URL(discoveryDoc.jwks_uri);
	} catch (error) {
		logger.error('Failed to retrieve JWKS URI from discovery endpoint', error);
		throw error; // Re-throw to indicate failure
	}
}

function extractClaimValue(payload: Record<string, unknown>, claimName: string): string[] {
	const value = payload[claimName];
	if (Array.isArray(value)) {
		return value.filter((item) => typeof item === 'string');
	}
	if (typeof value === 'string') {
		return value.split(' ').filter((item) => item.length > 0);
	}
	return [];
}

export async function getAuthUser(cookies: {
	get(name: string): string | undefined;
}): Promise<UserClaims | null> {
	const idToken = cookies.get('id_token');
	const accessToken = cookies.get('access_token');

	if (!idToken || !accessToken) {
		return null;
	}

	try {
		// 1. Get the JWKS URI from the discovery endpoint
		const jwksUri = await getJwksUri(env.PUBLIC_OIDC_ISSUER);

		// 2. Create a remote JWK set
		const JWKS = jose.createRemoteJWKSet(jwksUri);

		// 3. Verify the id_token
		const { payload: idPayload } = await jose.jwtVerify(idToken, JWKS, {
			issuer: env.PUBLIC_OIDC_ISSUER, // Corrected: Verify the issuer
			audience: env.PUBLIC_OIDC_CLIENT_ID // Verify the audience
			// Clock tolerance defaults to 60 seconds, adjust if needed
		});

		// 4. Verify the access_token
		const { payload: accessPayload } = await jose.jwtVerify(accessToken, JWKS, {
			issuer: env.PUBLIC_OIDC_ISSUER,
			audience: env.PUBLIC_OIDC_CLIENT_ID
		});

		// 5. Construct the UserClaims object from the validated payloads
		// Ensure payload properties used exist and have expected types (jose helps here)
		const groupsClaim = env.PUBLIC_OIDC_GROUPS_CLAIM || 'groups';
		const rolesClaim = env.PUBLIC_OIDC_ROLES_CLAIM || 'roles';

		const user: UserClaims = {
			sub: idPayload.sub ?? '', // Use ?? for required fields if payload types are loose
			email: (idPayload.email || idPayload.preferred_username || idPayload.sub)?.toString() ?? '',
			name: (idPayload.name || idPayload.given_name)?.toString(),
			roles: extractClaimValue(accessPayload, rolesClaim),
			groups: extractClaimValue(accessPayload, groupsClaim),
			id_token: idToken,
			access_token: accessToken
		};

		// Validate essential claims are present after extraction
		if (!user.sub || !user.email) {
			logger.error('Validated tokens missing required claims (sub or email)', {
				idPayload,
				accessPayload
			});
			return null;
		}

		logger.debug(
			{ userId: user.sub, groups: user.groups, roles: user.roles },
			'Tokens verified successfully'
		);
		return user;
	} catch (error: unknown) {
		if (error instanceof jose.errors.JWTExpired) {
			logger.info('Auth token validation failed: Token expired');
		} else if (
			error instanceof jose.errors.JOSEError || // Catches various jose errors (signature, claims, etc.)
			error instanceof Error
		) {
			logger.warn(`Auth token validation failed: ${error.message}`, { error });
		} else {
			logger.error('An unexpected error occurred during token validation', { error });
		}
		return null;
	}
}
