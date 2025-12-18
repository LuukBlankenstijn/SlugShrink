<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.ico';
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import type { LayoutData } from './$types';
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import type { Snippet } from 'svelte';
	import { setContext } from 'svelte';
	import Modal from '$lib/Modal.svelte';
	import { MODAL_CONTEXT, type ModalControls } from '$lib/modal-context';
	import type { Domain, FullRedirect } from '../gen/api/v1/api_pb';
	import RedirectForm from '$lib/RedirectForm.svelte';
	import DomainForm from '$lib/DomainForm.svelte';

	let open = $state(false);
	let selectedDomain = $state<Domain | null>(null);
	let selectedRedirect = $state<FullRedirect | null>(null);
	const active = $derived(() =>
		page.url.pathname.startsWith('/domains') ? 'domains' : 'redirects'
	);
	const mode = $derived(() =>
		selectedDomain ? 'domains' : selectedRedirect ? 'redirects' : active()
	);
	const modalTitle = $derived(() => {
		if (mode() === 'domains') return selectedDomain ? 'Edit domain' : 'New domain';
		return selectedRedirect ? 'Edit redirect' : 'New redirect';
	});
	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	function openDomainModal(domain?: Domain) {
		selectedDomain = domain ?? null;
		selectedRedirect = null;
		open = true;
	}

	function openRedirectModal(redirect?: FullRedirect) {
		selectedRedirect = redirect ?? null;
		selectedDomain = null;
		open = true;
	}

	function closeModal() {
		open = false;
		selectedDomain = null;
		selectedRedirect = null;
	}

	setContext<ModalControls>(MODAL_CONTEXT, {
		openDomainModal,
		openRedirectModal,
		closeModal
	});

	$effect(() => {
		if (!open) {
			selectedDomain = null;
			selectedRedirect = null;
		}
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<QueryClientProvider client={data.queryClient}>
	<div class="min-h-svh bg-slate-950 text-slate-100">
		<div class="mx-auto flex max-w-6xl flex-col gap-6 px-4 py-10">
			<header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
				<div class="space-y-2">
					<div class="flex gap-4">
						<a
							href={resolve('/redirects')}
							class="relative text-2xl font-semibold tracking-tight text-white/60 hover:text-white"
							class:active={active() === 'redirects'}
						>
							Redirects
							{#if active() === 'redirects'}
								<span
									class="absolute -bottom-2 left-0 h-px w-full bg-linear-to-r from-cyan-400/70 via-blue-500/70 to-violet-500/70"
								></span>
							{/if}
						</a>

						<a
							href={resolve('/domains')}
							class="relative text-2xl font-semibold tracking-tight text-white/60 hover:text-white"
							class:active={active() === 'domains'}
						>
							Domains
							{#if active() === 'domains'}
								<span
									class="absolute -bottom-2 left-0 h-px w-full bg-linear-to-r from-cyan-400/70 via-blue-500/70 to-violet-500/70"
								></span>
							{/if}
						</a>
					</div>
				</div>

				{#if active() === 'domains' || active() === 'redirects'}
					<div class="flex items-center gap-3">
						<button
							type="button"
							class="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold text-white hover:bg-white/10 focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/40"
							onclick={() => (active() === 'domains' ? openDomainModal() : openRedirectModal())}
						>
							New
						</button>
					</div>
				{/if}
			</header>
			<Modal bind:open title={modalTitle()}>
				{#if active() === 'domains'}
					<DomainForm current={selectedDomain ?? undefined} onClose={closeModal}></DomainForm>
				{:else if active() === 'redirects'}
					<RedirectForm current={selectedRedirect ?? undefined} onClose={closeModal}></RedirectForm>
				{/if}
			</Modal>
			{@render children()}
		</div>
	</div>
</QueryClientProvider>
