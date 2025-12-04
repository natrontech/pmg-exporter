// src/app.d.ts
declare global {
	namespace App {
		interface Locals {
			user?: UserClaims;
		}
	}
}

declare module '$houdini' {
	interface Session {
		user?: UserClaims;
	}
}

interface UserClaims {
	sub: string;
	email: string;
	name?: string;
	roles: string[];
	groups: string[];
	id_token: string;
	access_token: string;
	exp?: number;
}

export {};
