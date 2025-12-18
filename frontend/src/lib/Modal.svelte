<script lang="ts">
	import { tick, type Snippet } from 'svelte';

	let {
		open = $bindable(),
		title = 'Modal',
		children
	} = $props<{
		open: boolean;
		title?: string;
		children: Snippet;
	}>();

	let closeButton = $state<HTMLButtonElement | null>(null);

	function close() {
		open = false;
	}

	function onKeydown(event: KeyboardEvent) {
		if (event.key === 'Escape') close();
	}

	$effect(() => {
		if (!open) return;
		void tick().then(() => closeButton?.focus());
	});
</script>

<svelte:window
	onkeydown={(event) => {
		if (open) onKeydown(event);
	}}
/>

{#if open}
	<div
		class="fixed inset-0 z-50"
		role="presentation"
		onclick={(event) => {
			const target = event.target as HTMLElement;
			if (event.currentTarget === target || target.dataset.overlay === 'true') close();
		}}
	>
		<div class="absolute inset-0 bg-black/60 backdrop-blur-[2px]" data-overlay="true"></div>

		<div class="absolute inset-0 flex items-center justify-center p-4">
			<div
				class="w-full max-w-xl overflow-hidden rounded-2xl border border-white/10 bg-slate-950 text-slate-100 shadow-2xl shadow-black/30"
				role="dialog"
				aria-modal="true"
				aria-label={title}
				tabindex="-1"
			>
				<div class="flex items-center justify-between gap-4 border-b border-white/10 px-5 py-4">
					<div class="min-w-0">
						<p class="truncate text-sm font-semibold">{title}</p>
						<div
							class="mt-2 h-px w-24 bg-linear-to-r from-cyan-400/60 via-blue-500/60 to-violet-500/60"
						></div>
					</div>

					<button
						bind:this={closeButton}
						type="button"
						class="inline-flex h-9 w-9 items-center justify-center rounded-xl border border-white/10 bg-white/5 text-white/80 hover:bg-white/10 hover:text-white focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/40"
						aria-label="Close modal"
						onclick={close}
					>
						<span class="text-lg leading-none">×</span>
					</button>
				</div>

				<div class="px-5 py-4">
					{@render children()}
				</div>
			</div>
		</div>
	</div>
{/if}
