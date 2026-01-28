<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import type { Domain } from '../gen/api/v1/domain_pb';
	import { api } from './api.svelte';
	import Modal from '$lib/Modal.svelte';
	import { queryKeys } from '$lib/queryKeys';
	import { deleteDomainMutationOptions, saveDomainMutationOptions } from '$lib/mutationOptions';

	interface Props {
		readonly onClose: () => void;
		readonly current?: Domain;
	}

	const { current, onClose }: Props = $props();

	const qc = useQueryClient();
	let domain = $state('');
	let name = $state('');
	let error = $state<string | null>(null);
	let confirmDelete = $state(false);
	const isExisting = $derived(() => Boolean(current));
	let isViewMode = $state(Boolean(current));

	$effect(() => {
		domain = current?.domain ?? '';
		name = current?.name ?? '';
		isViewMode = Boolean(current);
	});

	const close = () => {
		onClose();
		domain = '';
		name = '';
		confirmDelete = false;
	};

	const cancelEdit = () => {
		if (!current) return;
		domain = current.domain ?? '';
		name = current.name ?? '';
		error = null;
		isViewMode = true;
	};

	const deleteDomain = createMutation(() => ({
		...deleteDomainMutationOptions(api),
		onSuccess: async () => {
			qc.invalidateQueries({ queryKey: queryKeys.domains() });
			close();
		},
		onError: (e) => {
			error = e instanceof Error ? e.message : 'Something went wrong';
		}
	}));

	const saveDomain = createMutation(() => ({
		...saveDomainMutationOptions(api, current),
		onSuccess: async () => {
			qc.invalidateQueries({ queryKey: queryKeys.domains() });
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
		if (isExisting() && isViewMode) {
			isViewMode = false;
			return;
		}
		saveDomain.mutate({ name, domain });
	}}
>
	<label class="block">
		<span class="text-sm text-slate-200">Name</span>
		{#if isViewMode}
			<div class="mt-1 rounded-lg bg-white/5 px-3 py-2 text-sm text-slate-200" title={name}>
				<span class="block truncate">{name || '—'}</span>
			</div>
		{:else}
			<input
				class="mt-1 w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
				bind:value={name}
				placeholder="example"
				autocomplete="off"
			/>
		{/if}
	</label>

	<label class="block">
		<span class="text-sm text-slate-200">Domain</span>
		{#if isViewMode}
			<div class="mt-1 rounded-lg bg-white/5 px-3 py-2 text-sm text-slate-200" title={domain}>
				<span class="block truncate">{domain || '—'}</span>
			</div>
		{:else}
			<input
				class="mt-1 w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
				bind:value={domain}
				placeholder="example.com"
				autocomplete="off"
			/>
		{/if}
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
		<div class="flex items-center gap-3">
			<button
				type="submit"
				class={[
					'inline-flex items-center justify-center rounded-xl border px-4 py-2 text-sm font-semibold transition focus-visible:ring-2 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-60',
					isViewMode
						? 'border-white/10 bg-white/10 text-white hover:bg-white/15 focus-visible:ring-cyan-400/40'
						: 'border-emerald-400/40 bg-emerald-500/20 text-emerald-50 hover:bg-emerald-500/30 focus-visible:ring-emerald-400/40'
				].join(' ')}
				disabled={saveDomain.isPending}
			>
				{#if saveDomain.isPending}
					Saving…
				{:else if isExisting() && isViewMode}
					Edit
				{:else}
					Save
				{/if}
			</button>

			{#if isExisting() && !isViewMode}
				<button
					type="button"
					class="inline-flex items-center justify-center rounded-xl border border-white/15 bg-white/5 px-4 py-2 text-sm font-semibold text-white hover:bg-white/10 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
					onclick={(e) => {
						e.preventDefault();
						cancelEdit();
					}}
				>
					Cancel
				</button>
			{/if}
		</div>

		{#if current && !isViewMode}
			<div class="flex items-center gap-3">
				<button
					type="button"
					onclick={(e) => {
						e.preventDefault();
						confirmDelete = true;
					}}
					class="inline-flex items-center justify-center rounded-xl border border-red-500/40 bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-100 hover:bg-red-500/20 focus-visible:ring-2 focus-visible:ring-red-400/40 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-60"
					disabled={deleteDomain.isPending}
				>
					Delete
				</button>

				<Modal bind:open={confirmDelete} title="Delete domain?">
					<div class="space-y-4">
						<p class="text-sm text-slate-200">
							This will remove <span class="font-semibold">{current.domain}</span> and all redirects using
							it. You can’t undo this action.
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
								onclick={() => deleteDomain.mutate({ id: current.id })}
								disabled={deleteDomain.isPending}
							>
								{deleteDomain.isPending ? 'Deleting…' : 'Delete'}
							</button>
						</div>
					</div>
				</Modal>
			</div>
		{/if}
	</div>
</form>
