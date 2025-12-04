// houdini.config.js
/** @type {import('houdini').ConfigFile} */
const config = {
	watchSchema: {
		url: 'http://localhost:8080/graphql',
		interval: 10000,
		timeout: 5000,
		headers: () => ({
			Accept: 'application/json'
		})
	},
	plugins: {
		'houdini-svelte': {
			client: './src/client.ts', // Explicit client path
			forceRunesMode: true
		}
	}
};
export default config;
