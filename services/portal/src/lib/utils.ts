/**
 * Simple string hash function (basic implementation)
 */
function simpleHash(str: string): number {
	let hash = 0;
	for (let i = 0; i < str.length; i++) {
		const char = str.charCodeAt(i);
		hash = (hash << 5) - hash + char;
		hash |= 0; // Convert to 32bit integer
	}
	return Math.abs(hash);
}

/**
 * Generates a CSS linear gradient string based on a user identifier
 */
export function generateAvatarGradient(identifier: string): string {
	if (!identifier) {
		// Default gradient for logged-out or missing ID
		return 'linear-gradient(45deg, #ccc, #eee)';
	}

	const hash = simpleHash(identifier);

	// Use hash to determine HSL values for two colors
	const h1 = hash % 360;
	const s1 = 50 + (hash % 30);
	const l1 = 40 + (hash % 20);

	const h2 = (h1 + 120 + (hash % 60)) % 360; // Ensure a decent color shift
	const s2 = 50 + ((hash >> 8) % 30);
	const l2 = 50 + ((hash >> 8) % 20);

	const color1 = `hsl(${h1}, ${s1}%, ${l1}%)`;
	const color2 = `hsl(${h2}, ${s2}%, ${l2}%)`;

	return `linear-gradient(45deg, ${color1}, ${color2})`;
}
