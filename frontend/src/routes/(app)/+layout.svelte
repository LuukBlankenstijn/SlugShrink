<script lang="ts">
	import { resolve } from '$app/paths';
	import { page } from '$app/state';
	import { goto } from '$app/navigation';
	import type { Snippet } from 'svelte';
	import { setContext } from 'svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { authStatusQueryOptions } from '$lib/queries/authStatus';
	import Modal from '$lib/Modal.svelte';
	import { MODAL_CONTEXT, type ModalControls } from '$lib/modal-context';
	import RedirectForm from '$lib/RedirectForm.svelte';
	import DomainForm from '$lib/DomainForm.svelte';
	import type { Domain } from '../../gen/api/v1/domain_pb';
	import type { FullRedirect } from '../../gen/api/v1/redirect_pb';
	import { UserPermission } from '../../gen/api/v1/auth_pb';

	let open = $state(false);
	let selectedDomain = $state<Domain | null>(null);
	let selectedRedirect = $state<FullRedirect | null>(null);
	const active = $derived(() => {
		if (page.url.pathname.startsWith('/domains')) return 'domains';
		if (page.url.pathname.startsWith('/settings')) return 'settings';
		return 'redirects';
	});
	const mode = $derived(() =>
		selectedDomain ? 'domains' : selectedRedirect ? 'redirects' : active()
	);
	const modalTitle = $derived(() => {
		if (mode() === 'domains') return selectedDomain ? 'Edit domain' : 'New domain';
		if (mode() === 'redirects') return selectedRedirect ? 'Edit redirect' : 'New redirect';
		return '';
	});
	let { children }: { children: Snippet } = $props();

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

	const authStatus = createQuery(authStatusQueryOptions);

	$effect(() => {
		if (authStatus.data === false) {
			goto(resolve('/login'));
		}
	});

	$effect(() => {
		if (!open) {
			selectedDomain = null;
			selectedRedirect = null;
		}
	});
</script>

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

					{#if authStatus.data && [UserPermission.PERMISSION_ADMIN, UserPermission.PERMISSION_SUPERUSER].includes(authStatus.data.permission)}
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
					{/if}

					{#if authStatus.data && authStatus.data.permission === UserPermission.PERMISSION_ADMIN}
						<a
							href={resolve('/settings')}
							class="relative text-2xl font-semibold tracking-tight text-white/60 hover:text-white"
							class:active={active() === 'settings'}
						>
							Settings
							{#if active() === 'settings'}
								<span
									class="absolute -bottom-2 left-0 h-px w-full bg-linear-to-r from-cyan-400/70 via-blue-500/70 to-violet-500/70"
								></span>
							{/if}
						</a>
					{/if}
				</div>
			</div>

			<div class="flex min-h-[42px] items-center gap-3">
				{#if active() === 'domains' || active() === 'redirects'}
					<button
						type="button"
						class="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold text-white hover:bg-white/10 focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/40"
						onclick={() => (active() === 'domains' ? openDomainModal() : openRedirectModal())}
					>
						New
					</button>
				{/if}
			</div>
		</header>
		{#if active() !== 'settings'}
			<Modal bind:open title={modalTitle()}>
				{#if active() === 'domains'}
					<DomainForm current={selectedDomain ?? undefined} onClose={closeModal}></DomainForm>
				{:else if active() === 'redirects'}
					<RedirectForm current={selectedRedirect ?? undefined} onClose={closeModal}></RedirectForm>
				{/if}
			</Modal>
		{/if}
		{@render children()}
	</div>
</div>
