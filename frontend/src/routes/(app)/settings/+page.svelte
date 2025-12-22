<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import BasicAuthSetting from './basicAuthSetting.svelte';
	import { getMethod, methods } from './authMethods';
	import AuthlessSettings from './authlessSettings.svelte';
	import ProxyAuthSettings from './proxyAuthSettings.svelte';
	import { authStatusQueryOptions } from '$lib/queries/authStatus';
	import { authConfigQueryOptions } from '$lib/queryOptions';
	import { hasPermission } from '$lib/auth';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { UserPermission } from '../../../gen/api/v1/auth_pb';
	const authStatus = createQuery(authStatusQueryOptions);

	$effect(() => {
		if (authStatus.data && !hasPermission(authStatus.data, UserPermission.PERMISSION_ADMIN)) {
			goto(resolve('/redirects'));
		}
	});

	let saveTrigger = $state(0);

	const query = createQuery(() => authConfigQueryOptions());

	let userSelection = $state<string | null>(null);
	let selection = $derived(userSelection ?? query.data?.config.case ?? '');
	$effect(() => {
		if (query.data?.config.case && !userSelection) {
			userSelection = query.data.config.case;
		}
	});
</script>

<svelte:head>
	<title>Settings · gewi.sh</title>
</svelte:head>

<section class="space-y-8">
	{#if query.data}
		<div
			class="overflow-hidden rounded-2xl border border-white/10 bg-white/5 shadow-[0_18px_80px_rgba(0,0,0,0.35)]"
		>
			<div
				class="flex flex-col gap-4 border-b border-white/10 px-6 py-4 sm:flex-row sm:items-center sm:justify-between"
			>
				<div class="space-y-1">
					<p class="text-sm font-semibold text-slate-100">Authentication method</p>
					<p class="text-xs text-slate-400">Select how users sign in.</p>
				</div>

				<div class="flex w-full max-w-md items-center gap-3">
					<div class="relative w-full max-w-xs">
						<select
							class="w-full appearance-none rounded-xl border border-white/15 bg-slate-950/70 px-4 py-2 pr-10 text-sm text-slate-50 transition outline-none hover:border-cyan-400/40 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
							aria-label="Authentication method"
							value={selection}
							onchange={(e) => {
								userSelection = (e.target as HTMLSelectElement).value;
							}}
						>
							{#each Object.entries(methods) as [key, method] (key)}
								<option value={key}>{method}</option>
							{/each}
						</select>
						<span
							class="pointer-events-none absolute inset-y-0 right-3 flex items-center text-slate-400"
							aria-hidden="true"
						>
							⌄
						</span>
					</div>

					<div
						class="flex items-center gap-2 rounded-full border border-cyan-400/25 bg-cyan-400/10 px-3 py-1 text-xs font-semibold text-cyan-100"
					>
						<span class="h-2 w-2 rounded-full bg-cyan-300 shadow-[0_0_12px_rgba(34,211,238,0.65)]"
						></span>
						<span class="whitespace-nowrap">Active: {getMethod(query.data.config.case)}</span>
					</div>
				</div>
			</div>

			<div class="grid gap-6 px-6 py-6 sm:grid-cols-[1fr_minmax(320px,1.4fr)]">
				<div class="space-y-2">
					<p class="text-sm font-semibold text-slate-100">Method details</p>
					<p class="text-sm text-slate-300">Fields update based on the selection.</p>
				</div>

				{#if userSelection === 'basicAuth'}
					<BasicAuthSetting {saveTrigger}></BasicAuthSetting>
				{:else if userSelection === 'proxyAuth'}
					{#if query.data.config.case === 'proxyAuth'}
						<ProxyAuthSettings {saveTrigger} currentConfig={query.data.config.value}
						></ProxyAuthSettings>
					{:else}
						<ProxyAuthSettings {saveTrigger}></ProxyAuthSettings>
					{/if}
				{:else}
					<AuthlessSettings {saveTrigger}></AuthlessSettings>
				{/if}
			</div>

			<div
				class="flex items-center justify-end gap-3 border-t border-white/10 bg-slate-950/30 px-6 py-4"
			>
				<button
					type="button"
					class="inline-flex items-center justify-center rounded-xl border border-white/10 px-4 py-2 text-sm font-semibold text-slate-100 hover:bg-white/10 focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/40"
				>
					Cancel
				</button>
				<button
					type="button"
					class="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/10 px-4 py-2 text-sm font-semibold text-white transition hover:bg-white/15 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
					onclick={() => saveTrigger++}
				>
					Save changes
				</button>
			</div>
		</div>
	{/if}
</section>
