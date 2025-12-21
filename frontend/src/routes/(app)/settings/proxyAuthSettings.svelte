<script lang="ts">
	import { api } from '$lib/api.svelte';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { untrack } from 'svelte';
	import { authStatusQueryKey } from '$lib/queries/authStatus';
	import { queryKeys } from '$lib/queryKeys';
	import type { ProxyAuth } from '../../../gen/api/v1/auth_pb';

	const { saveTrigger, currentConfig }: { saveTrigger: number; currentConfig?: ProxyAuth } =
		$props();
	const queryclient = useQueryClient();
	let lastProcessedTrigger = saveTrigger;

	let groupHeader = $state(currentConfig?.groupHeader ?? '');
	let userIdHeader = $state(currentConfig?.userIdHeader ?? '');
	let superUserGroups = $state(currentConfig?.superUserGroups ?? []);
	let adminGroups = $state(currentConfig?.adminGroups ?? []);
	let superUserGroupInput = $state('');
	let adminGroupInput = $state('');
	let groupSeparator = $state(currentConfig?.groupsSeparator ?? '');
	let validationError = $state<string | null>(null);
	let serverError = $state<string | null>(null);

	const errorMessage = $derived(validationError ?? serverError);

	const mutation = createMutation(() => ({
		mutationFn: () =>
			api.setAuthConfig({
				config: {
					case: 'proxyAuth',
					value: {
						groupHeader,
						userIdHeader,
						superUserGroups,
						adminGroups,
						groupsSeparator: groupSeparator ?? undefined
					}
				}
			}),
		onError: (err: unknown) => {
			serverError = err instanceof Error ? err.message : 'Something went wrong';
		},
		onSuccess: () => {
			validationError = null;
			serverError = null;
			queryclient.invalidateQueries({ queryKey: queryKeys.authConfig() });
			queryclient.invalidateQueries({ queryKey: authStatusQueryKey });
		},
		retry: false
	}));

	function commitPendingInputs() {
		if (superUserGroupInput.trim()) addSuperUserGroup(superUserGroupInput);
		if (adminGroupInput.trim()) addAdminGroup(adminGroupInput);
	}

	function validate() {
		commitPendingInputs();
		if (!groupHeader?.trim()) return 'Group header is required';
		if (!userIdHeader?.trim()) return 'User ID header is required';
		if (!adminGroups.length) return 'At least one admin group is required';
		// groupSeparator is optional
		return null;
	}

	function addSuperUserGroup(value: string) {
		const trimmed = value.trim();
		if (!trimmed || superUserGroups.includes(trimmed)) {
			superUserGroupInput = '';
			return;
		}
		superUserGroups = [...superUserGroups, trimmed];
		superUserGroupInput = '';
	}

	function addAdminGroup(value: string) {
		const trimmed = value.trim();
		if (!trimmed || adminGroups.includes(trimmed)) {
			adminGroupInput = '';
			return;
		}
		adminGroups = [...adminGroups, trimmed];
		adminGroupInput = '';
	}

	function removeSuperUserGroup(value: string) {
		superUserGroups = superUserGroups.filter((v) => v !== value);
	}

	function removeAdminGroup(value: string) {
		adminGroups = adminGroups.filter((v) => v !== value);
	}

	function handleSuperUserKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addSuperUserGroup(superUserGroupInput);
		}
	}

	function handleAdminKeydown(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			event.preventDefault();
			addAdminGroup(adminGroupInput);
		}
	}

	$effect(() => {
		if (saveTrigger > lastProcessedTrigger && !mutation.isPending) {
			lastProcessedTrigger = saveTrigger;
			untrack(() => {
				const validation = validate();
				if (validation) {
					validationError = validation;
					return;
				}
				mutation.mutate();
			});
		}
	});

	$effect(() => {
		validationError = null;
		serverError = null;
	});
</script>

<div class="space-y-4 rounded-xl border border-white/10 bg-slate-950/40 p-4">
	<h4 class="whitespace-nowrap">Update Settings:</h4>
	<div>
		<label for="groupHeader-input" class="flex items-center gap-2 text-sm text-slate-200">
			Group header
			<span
				class="inline-flex h-4 w-4 items-center justify-center rounded-full border border-slate-500/60 text-xs text-slate-200"
				title="Header containing group memberships">?</span
			>
		</label>
		<input
			bind:value={groupHeader}
			id="groupHeader-input"
			class="w-full rounded-lg border border-white/10 bg-slate-900/70 px-3 py-2 pr-12 text-sm text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
			placeholder="e.g. X-Proxy-Groups"
		/>
	</div>
	<div>
		<label for="userIdHeader-input" class="flex items-center gap-2 text-sm text-slate-200">
			User ID header
			<span
				class="inline-flex h-4 w-4 items-center justify-center rounded-full border border-slate-500/60 text-xs text-slate-200"
				title="Header containing the unique user identifier">?</span
			>
		</label>
		<input
			bind:value={userIdHeader}
			id="userIdHeader-input"
			class="w-full rounded-lg border border-white/10 bg-slate-900/70 px-3 py-2 pr-12 text-sm text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
			placeholder="e.g. X-Proxy-User"
		/>
	</div>
	<div>
		<label for="superUserGroup-input" class="flex items-center gap-2 text-sm text-slate-200">
			Super user group
			<span
				class="inline-flex h-4 w-4 items-center justify-center rounded-full border border-slate-500/60 text-xs text-slate-200"
				title="Groups that grant super user rights; add multiple">?</span
			>
		</label>
		<div class="space-y-2">
			<input
				bind:value={superUserGroupInput}
				onkeydown={handleSuperUserKeydown}
				id="superUserGroup-input"
				class="w-full rounded-lg border border-white/10 bg-slate-900/70 px-3 py-2 pr-12 text-sm text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
				placeholder="Type a group name and press Enter"
			/>
			<div class="flex flex-wrap gap-2">
				{#each superUserGroups as group (group)}
					<span
						class="inline-flex items-center gap-2 rounded-full border border-cyan-400/40 bg-cyan-500/10 px-3 py-1 text-xs text-cyan-100"
					>
						{group}
						<button
							type="button"
							class="text-slate-300 hover:text-white"
							onclick={() => removeSuperUserGroup(group)}
						>
							x
						</button>
					</span>
				{/each}
			</div>
		</div>
	</div>
	<div>
		<label for="adminGroup-input" class="flex items-center gap-2 text-sm text-slate-200">
			Admin group
			<span
				class="inline-flex h-4 w-4 items-center justify-center rounded-full border border-slate-500/60 text-xs text-slate-200"
				title="Groups that grant admin rights; add multiple">?</span
			>
		</label>
		<div class="space-y-2">
			<input
				bind:value={adminGroupInput}
				onkeydown={handleAdminKeydown}
				id="adminGroup-input"
				class="w-full rounded-lg border border-white/10 bg-slate-900/70 px-3 py-2 pr-12 text-sm text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
				placeholder="Type a group name and press Enter"
			/>
			<div class="flex flex-wrap gap-2">
				{#each adminGroups as group (group)}
					<span
						class="inline-flex items-center gap-2 rounded-full border border-cyan-400/40 bg-cyan-500/10 px-3 py-1 text-xs text-cyan-100"
					>
						{group}
						<button
							type="button"
							class="text-slate-300 hover:text-white"
							onclick={() => removeAdminGroup(group)}
						>
							x
						</button>
					</span>
				{/each}
			</div>
		</div>
	</div>
	<div>
		<label for="groupSeparator-input" class="flex items-center gap-2 text-sm text-slate-200">
			Group separator
			<span
				class="inline-flex h-4 w-4 items-center justify-center rounded-full border border-slate-500/60 text-xs text-slate-200"
				title="Separator used between groups in the group header; defaults to '|' when empty"
				>?</span
			>
		</label>
		<input
			bind:value={groupSeparator}
			id="groupSeparator-input"
			class="w-full rounded-lg border border-white/10 bg-slate-900/70 px-3 py-2 pr-12 text-sm text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
			placeholder="Defaults to |"
		/>
	</div>

	<div class="space-y-1"></div>

	{#if errorMessage}
		<p class="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">
			{errorMessage}
		</p>
	{/if}
</div>
