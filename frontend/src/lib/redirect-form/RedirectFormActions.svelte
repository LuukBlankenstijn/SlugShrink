<script lang="ts">
	import Modal from '$lib/Modal.svelte';
	import type { FullRedirect } from '../../gen/api/v1/redirect_pb';

	export let isExisting = false;
	export let isViewMode = true;
	export let canEdit = true;
	export let isSaving = false;
	export let canCreate = true;
	export let onCancel: () => void = () => {};
	export let onDelete: (id: string) => void = () => {};
	export let isDeleting = false;
	export let current: FullRedirect | undefined;

	let confirmDelete = false;

	const startDelete = (e: Event) => {
		e.preventDefault();
		confirmDelete = true;
	};
</script>

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
			disabled={isSaving || (!isExisting && !canCreate) || (isExisting && isViewMode && !canEdit)}
			title={isExisting && isViewMode && !canEdit ? 'You are not the owner of this redirect' : ''}
		>
			{#if isSaving}
				Saving…
			{:else if isExisting && isViewMode}
				Edit
			{:else if !isExisting && !canCreate}
				Add a domain first
			{:else}
				Save
			{/if}
		</button>

		{#if isExisting && !isViewMode}
			<button
				type="button"
				class="inline-flex items-center justify-center rounded-xl border border-white/15 bg-white/5 px-4 py-2 text-sm font-semibold text-white hover:bg-white/10 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
				onclick={(e) => {
					e.preventDefault();
					onCancel();
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
				onclick={startDelete}
				class="inline-flex items-center justify-center rounded-xl border border-red-500/40 bg-red-500/10 px-4 py-2 text-sm font-semibold text-red-100 hover:bg-red-500/20 focus-visible:ring-2 focus-visible:ring-red-400/40 focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-60"
				disabled={isDeleting}
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
							onclick={() => current && onDelete(current.id)}
							disabled={isDeleting}
						>
							{isDeleting ? 'Deleting…' : 'Delete'}
						</button>
					</div>
				</div>
			</Modal>
		</div>
	{/if}
</div>
