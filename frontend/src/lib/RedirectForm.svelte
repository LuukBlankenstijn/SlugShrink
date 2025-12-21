<script lang="ts">
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { api } from './api.svelte';
	import Modal from '$lib/Modal.svelte';
	import { queryKeys } from '$lib/queryKeys';
	import type { FullRedirect } from '../gen/api/v1/redirect_pb';

	interface Props {
		readonly onClose: () => void;
		readonly current?: FullRedirect;
	}

	const { onClose, current }: Props = $props();

	const qc = useQueryClient();
	const domainsQuery = createQuery(() => ({
		queryKey: queryKeys.domains(),
		queryFn: () => api.getDomains({})
	}));

	const domains = $derived(() => domainsQuery.data?.data ?? []);

	let domainId = $state('');
	let path = $state('');
	let targetUrl = $state('');
	let active = $state(true);
	let error = $state<string | null>(null);
	let confirmDelete = $state(false);

	const isEditing = $derived(() => Boolean(current));

	const close = () => {
		onClose();
		domainId = '';
		path = '';
		targetUrl = '';
		active = true;
		error = null;
		confirmDelete = false;
	};

	$effect(() => {
		if (current) {
			domainId = current.domainId;
			path = current.path;
			targetUrl = current.targetUrl;
			active = current.active;
		} else {
			domainId = '';
			path = '';
			targetUrl = '';
			active = true;
		}
		error = null;
	});

	$effect(() => {
		if (!domainId && domains().length > 0) {
			domainId = domains()[0].id;
		}
	});

	const createRedirect = createMutation(() => ({
		mutationFn: ({
			domainId,
			path,
			targetUrl,
			active
		}: {
			domainId: string;
			path: string;
			targetUrl: string;
			active: boolean;
		}) =>
			isEditing()
				? api.putRedirect({ id: current?.id ?? '', domainId, path, targetUrl, active })
				: api.createRedirect({ domainId, path, targetUrl, active }),
		onSuccess: async () => {
			qc.invalidateQueries({ queryKey: queryKeys.redirects() });
			close();
		},
		onError: (e) => {
			error = e instanceof Error ? e.message : 'Something went wrong';
		}
	}));

	const deleteRedirect = createMutation(() => ({
		mutationFn: ({ id }: { id: string }) => api.deleteRedirect({ id }),
		onSuccess: async () => {
			qc.invalidateQueries({ queryKey: queryKeys.redirects() });
			close();
		},
		onError: (e) => {
			error = e instanceof Error ? e.message : 'Something went wrong';
		}
	}));
</script>

<form
	class="space-y-4"
	onsubmit={(e) => {
		e.preventDefault();
		createRedirect.mutate({ domainId, path, targetUrl, active });
	}}
>
	<label class="block space-y-1.5">
		<span class="text-sm text-slate-200">Domain</span>
		<select
			class="w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20 disabled:opacity-60"
			bind:value={domainId}
			disabled={domainsQuery.isPending || domains().length === 0}
		>
			{#if domains().length === 0}
				<option value="">No domains available</option>
			{:else}
				{#each domains() as domain (domain.id)}
					<option value={domain.id}>{domain.domain}</option>
				{/each}
			{/if}
		</select>
	</label>

	<label class="block space-y-1.5">
		<span class="text-sm text-slate-200">Path</span>
		<input
			class="w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
			bind:value={path}
			placeholder="/promo"
			autocomplete="off"
		/>
	</label>

	<label class="block space-y-1.5">
		<span class="text-sm text-slate-200">Target URL</span>
		<input
			class="w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
			bind:value={targetUrl}
			placeholder="https://example.com/landing"
			autocomplete="off"
		/>
	</label>

	<label
		class="flex w-fit items-center gap-2 rounded-xl border border-white/10 bg-slate-950/30 px-3 py-2 text-sm text-slate-200 select-none"
	>
		<input type="checkbox" class="h-4 w-4 accent-cyan-400" bind:checked={active} />
		<span>Active</span>
	</label>

	{#if error}
		<p class="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">
			{error}
		</p>
	{/if}

	<div
		class={['mt-2 flex items-center gap-3', current ? 'justify-between' : 'justify-start'].join(
			' '
		)}
	>
		<button
			type="submit"
			class="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/10 px-4 py-2 text-sm font-semibold text-white hover:bg-white/15 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-60"
			disabled={createRedirect.isPending || domains().length === 0}
		>
			{#if createRedirect.isPending}
				Saving…
			{:else if isEditing()}
				Update
			{:else if domains().length === 0}
				Add a domain first
			{:else}
				Save
			{/if}
		</button>

		{#if current}
			<div class="flex items-center gap-3">
				<button
					type="button"
					onclick={(e) => {
						e.preventDefault();
						confirmDelete = true;
					}}
					class="inline-flex items-center justify-center rounded-xl border border-red-500/40 bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-100 hover:bg-red-500/20 focus-visible:ring-2 focus-visible:ring-red-400/40 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-60"
					disabled={deleteRedirect.isPending}
				>
					Delete
				</button>

				<Modal bind:open={confirmDelete} title="Delete redirect?">
					<div class="space-y-4">
						<p class="text-sm text-slate-200">
							This will remove <span class="font-semibold">{current.domain}{current.path}</span> →
							<span class="font-semibold">{current.targetUrl}</span>. You can’t undo this action.
						</p>
						<div class="flex items-center justify-end gap-3">
							<button
								type="button"
								class="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/10 px-4 py-2 text-sm font-semibold text-white hover:bg-white/15 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
								onclick={() => (confirmDelete = false)}
							>
								Cancel
							</button>
							<button
								type="button"
								class="inline-flex items-center justify-center rounded-xl border border-red-500/40 bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-100 hover:bg-red-500/20 focus-visible:ring-2 focus-visible:ring-red-400/40 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-60"
								onclick={() => deleteRedirect.mutate({ id: current.id })}
								disabled={deleteRedirect.isPending}
							>
								{deleteRedirect.isPending ? 'Deleting…' : 'Delete'}
							</button>
						</div>
					</div>
				</Modal>
			</div>
		{/if}
	</div>
</form>
