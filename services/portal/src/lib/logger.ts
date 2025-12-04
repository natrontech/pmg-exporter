import { browser } from '$app/environment';
import pino from 'pino';

// --- Server-Side Logger Configuration ---
let serverLogger: pino.Logger | null = null;
if (!browser) {
	// Get log level from environment or default to 'info' (only on server)
	const level = process.env.LOG_LEVEL?.toLowerCase() || 'info';

	// Create a pino logger instance (only on server)
	serverLogger = pino({
		level,
		transport:
			process.env.NODE_ENV === 'test' // Disable pretty print in tests
				? undefined
				: {
						target: 'pino-pretty',
						options: {
							colorize: process.env.NODE_ENV !== 'production'
						}
					},
		base: undefined, // Don't add pid and hostname
		timestamp: pino.stdTimeFunctions.isoTime
	});
}

// --- Client-Side Logger Mock/Fallback ---
// Provides the same interface but uses console
const clientLogger = {
	trace: (...args: unknown[]) => console.trace(...args),
	debug: (...args: unknown[]) => console.debug(...args),
	info: (...args: unknown[]) => console.info(...args),
	warn: (...args: unknown[]) => console.warn(...args),
	error: (...args: unknown[]) => console.error(...args),
	fatal: (...args: unknown[]) => console.error('FATAL:', ...args) // No console.fatal
};

// --- Exported Logger ---
// Use serverLogger if available (server-side), otherwise use clientLogger (client-side)
export const logger = browser ? clientLogger : serverLogger!;

// Export log level constants for convenience (can be used anywhere)
export const LogLevel = {
	TRACE: 'trace',
	DEBUG: 'debug',
	INFO: 'info',
	WARN: 'warn',
	ERROR: 'error',
	FATAL: 'fatal'
};
