<script lang="ts">
	import { generateAvatarGradient } from '$lib/utils';
	import { toast } from 'svelte-sonner';
	import { LogOut, Home, User, Copy, Layers, Key } from 'lucide-svelte';
	import { logout } from '$lib/client/oidc';
	import type { LayoutData } from '../../routes/(app)/$types';
	import { page } from '$app/state';

	let { data } = $props<{ data: LayoutData }>();

	// Calculate avatar details
	const avatarGradient = generateAvatarGradient(data.user?.sub || '');
	const userInitial = data.user?.email ? data.user.email[0].toUpperCase() : '?';

	// Define navigation items
	const navItems = [
		{ href: '/overview', icon: Home, label: 'Overview' },
		{ href: '/namespaces', icon: Layers, label: 'Namespaces' },
		{ href: '/profile', icon: User, label: 'Profile' }
	];

	// Function to determine active state based on reactive $page
	function isActive(href: string): boolean {
		const path = page.url.pathname;
		if (href === '/overview') {
			return path === href || path.startsWith(href + '/');
		}
		return path.startsWith(href);
	}

	// Function to copy a token
	async function copyToken(tokenType: 'access' | 'id') {
		const token = tokenType === 'access' ? data.user?.access_token : data.user?.id_token;
		const tokenName = tokenType === 'access' ? 'Access Token' : 'ID Token';

		if (token && navigator.clipboard) {
			try {
				await navigator.clipboard.writeText(token);
				toast.success(`${tokenName} copied to clipboard`);
			} catch (err) {
				toast.error(`Failed to copy ${tokenName.toLowerCase()}: ${err}`);
			}
		} else if (!token) {
			toast.error(`${tokenName} not available in user data.`);
		} else {
			toast.error('Clipboard API not available or permission denied.');
		}
	}
</script>

<div class="flex h-full flex-col">
	<nav class="flex-1 overflow-y-auto px-2 py-4">
		<ul class="menu menu-md w-full space-y-1">
			{#each navItems as item (item.href)}
				<li>
					<a
						href={item.href}
						class="flex items-center gap-3 rounded-md {isActive(item.href)
							? 'bg-primary text-primary-content font-semibold'
							: ''}"
					>
						<item.icon class="h-5 w-5" />
						<span class="text-sm font-medium">{item.label}</span>
					</a>
				</li>
			{/each}
		</ul>
	</nav>

	<div class="border-base-300 border-t p-2">
		<div class="dropdown dropdown-top w-full">
			<div
				tabindex="0"
				role="button"
				class="btn btn-ghost w-full items-center justify-start gap-3 px-2"
			>
				<div class="avatar">
					<div class="w-8 rounded-full" style="background: {avatarGradient};">
						<span
							class="flex h-full w-full items-center justify-center text-sm font-semibold text-white"
						>
							{userInitial}
						</span>
					</div>
				</div>
				<div class="flex-1 overflow-hidden text-left">
					<p class="truncate text-sm font-medium">{data.user?.email || 'User'}</p>
				</div>
			</div>
			<ul
				class="dropdown-content menu menu-sm bg-base-100 rounded-box border-base-300 z-[999] mb-1 w-full border shadow-lg"
			>
				<li>
					<button onclick={() => copyToken('access')} class="flex items-center gap-2 text-sm">
						<Key class="h-4 w-4" />Copy Access Token
					</button>
				</li>
				<li>
					<button onclick={() => copyToken('id')} class="flex items-center gap-2 text-sm">
						<Copy class="h-4 w-4" />Copy ID Token
					</button>
				</li>
				<li>
					<button onclick={logout} class="flex items-center gap-2 text-sm">
						<LogOut class="h-4 w-4" />Logout
					</button>
				</li>
			</ul>
		</div>
	</div>
</div>
