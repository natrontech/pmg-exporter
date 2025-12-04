import { QueryStore } from '../runtime/stores/query'
import artifact from '$houdini/artifacts/Namespace'
import { initClient } from '$houdini/plugins/houdini-svelte/runtime/client'

export class NamespaceStore extends QueryStore {
	constructor() {
		super({
			artifact,
			storeName: "NamespaceStore",
			variables: false,
		})
	}
}

export async function load_Namespace(params) {
  await initClient()

	const store = new NamespaceStore()

	await store.fetch(params)

	return {
		Namespace: store,
	}
}
