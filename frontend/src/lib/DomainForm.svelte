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

	$effect(() => {
		domain = current?.domain ?? '';
		name = current?.name ?? '';
	});

	const close = () => {
		onClose();
		domain = '';
		name = '';
		confirmDelete = false;
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
		saveDomain.mutate({ name, domain });
	}}
>
	<label class="block">
		<span class="text-sm text-slate-200">Name</span>
		<input
			class="mt-1 w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
			bind:value={name}
			placeholder="example"
			autocomplete="off"
		/>
	</label>

	<label class="block">
		<span class="text-sm text-slate-200">Domain</span>
		<input
			class="mt-1 w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
			bind:value={domain}
			placeholder="example.com"
			autocomplete="off"
		/>
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
			disabled={saveDomain.isPending}
		>
			{saveDomain.isPending ? 'Saving…' : 'Save'}
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
