<script lang="ts">
	import { api } from '$lib/api.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { MODAL_CONTEXT, type ModalControls } from '$lib/modal-context';
	import { queryKeys } from '$lib/queryKeys';
	import { authStatusQueryOptions } from '$lib/queries/authStatus';
	import { UserPermission } from '../../../gen/api/v1/auth_pb';

	let page = $state(1);
	let pageSize = $state(10);
	const query = createQuery(() => ({
		queryKey: queryKeys.redirects(page, pageSize),
		queryFn: () =>
			api.getRedirects({
				page,
				pagesize: pageSize
			}),
		placeholderData: () => ({ data: [], total: 0 }),
		keepPreviousData: true
	}));

	const redirects = $derived(() => query.data?.data ?? []);
	const total = $derived(() => query.data?.total ?? 0);
	const totalPages = $derived(() => Math.max(1, Math.ceil(total() / pageSize)));
	const placeholderRows = $derived(() => Array.from({ length: Math.max(4, Math.min(8, pageSize)) }));
	const { openRedirectModal } = getContext<ModalControls>(MODAL_CONTEXT);

	const authStatus = createQuery(authStatusQueryOptions);

	const canEdit = (creator?: string) => {
		return (
			authStatus.data &&
			(authStatus.data.userId === creator ||
				!(authStatus.data.permission === UserPermission.PERMISSION_USER))
		);
	};
</script>

<svelte:head>
	<title>Redirects · gewi.sh</title>
</svelte:head>

<section
	class="overflow-hidden rounded-2xl border border-white/10 bg-white/3"
	aria-label="Redirects table"
>
	<div
		class="flex flex-col gap-3 border-b border-white/10 p-4 sm:flex-row sm:items-center sm:justify-between"
	>
		<p class="text-sm text-slate-300">
			<span class="font-semibold text-slate-100">{total()}</span>
			{total() === 1 ? 'redirect' : 'redirects'} found
		</p>
		<div class="hidden items-center gap-3 text-xs text-slate-300 sm:flex">
			<button
				type="button"
				class="rounded-lg border border-white/15 bg-white/5 px-3 py-1.5 font-semibold text-white transition enabled:hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40"
				onclick={() => (page = Math.max(1, page - 1))}
				disabled={page === 1 || query.isFetching}
			>
				Prev
			</button>
			<span class="text-slate-200">
				Page {page} / {totalPages()}
			</span>
			<button
				type="button"
				class="rounded-lg border border-white/15 bg-white/5 px-3 py-1.5 font-semibold text-white transition enabled:hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40"
				onclick={() => (page = Math.min(totalPages(), page + 1))}
				disabled={page >= totalPages() || query.isFetching}
			>
				Next
			</button>
		</div>
	</div>

	<div class="overflow-x-auto">
		<table class="min-w-full divide-y divide-white/10 text-sm">
			<thead class="bg-white/3">
				<tr class="text-left text-xs font-semibold tracking-wide text-slate-300 uppercase">
					<th class="px-4 py-3">Status</th>
					<th class="px-4 py-3">Domain</th>
					<th class="px-4 py-3">Path</th>
					<th class="px-4 py-3">Target</th>
					<th class="px-4 py-3"></th>
				</tr>
			</thead>
			<tbody class="divide-y divide-white/10">
				{#if query.isLoading || (query.isFetching && !redirects().length)}
					{#each placeholderRows() as _, idx}
						<tr class="animate-pulse">
							<td class="px-4 py-3">
								<div class="h-6 w-20 rounded-full bg-white/10"></div>
							</td>
							<td class="px-4 py-3">
								<div class="h-5 w-40 rounded bg-white/10"></div>
							</td>
							<td class="px-4 py-3">
								<div class="h-6 w-24 rounded bg-white/10"></div>
							</td>
							<td class="px-4 py-3">
								<div class="h-5 w-52 rounded bg-white/10"></div>
							</td>
							<td class="px-4 py-3 text-right">
								<div class="ml-auto h-8 w-16 rounded bg-white/10"></div>
							</td>
						</tr>
					{/each}
				{:else}
					{#each redirects() as r (r.id)}
						<tr class="group hover:bg-white/4">
							<td class="px-4 py-3">
								<span
									class={[
										'inline-flex items-center gap-2 rounded-full border px-2.5 py-1 text-xs font-semibold',
										r.active
											? 'border-emerald-400/30 bg-emerald-400/10 text-emerald-200'
											: 'border-slate-400/20 bg-white/5 text-slate-300'
									].join(' ')}
								>
									<span
										class={[
											'h-1.5 w-1.5 rounded-full',
											r.active ? 'bg-emerald-300' : 'bg-slate-400'
										].join(' ')}
									></span>
									{r.active ? 'Active' : 'Paused'}
								</span>
							</td>
							<td class="px-4 py-3">
								<span class="inline-flex max-w-[16rem] items-center gap-2">
									<span class="h-2 w-2 rounded-full bg-cyan-400/80"></span>
									<span class="truncate text-slate-200" title={r.domain}>
										{r.domain}
									</span>
								</span>
							</td>
							<td class="px-4 py-3">
								<span
									class="inline-flex rounded-lg border border-white/10 bg-slate-950/30 px-2 py-1 font-mono text-xs text-slate-200"
									title={r.path}
								>
									{r.path}
								</span>
							</td>
							<td class="px-4 py-3">
								<span class="block max-w-136 truncate text-slate-200" title={r.targetUrl}>
									{r.targetUrl}
								</span>
							</td>
							<td class="px-4 py-3 text-right">
								<button
									type="button"
									class="inline-flex items-center justify-center rounded-lg border border-white/15 bg-white/8 px-3 py-1.5 text-xs font-semibold text-white transition focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none enabled:border-white/25 enabled:hover:bg-white/12 disabled:cursor-not-allowed disabled:opacity-50"
									onclick={() => openRedirectModal(r)}
									title={canEdit(r.creator) ? '' : 'You are not the owner of this redirect'}
									disabled={!canEdit(r.creator)}
								>
									Edit
								</button>
							</td>
						</tr>
					{/each}
				{/if}
			</tbody>
		</table>
	</div>
	<div class="flex flex-col gap-3 p-2 sm:flex-row sm:items-center sm:justify-between">
		<label class="flex items-center gap-2 text-xs text-slate-300">
			Page size
			<select
				class="rounded-lg border border-white/15 bg-slate-950/60 px-3 py-2 text-sm text-slate-100 outline-none focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
				value={pageSize}
				onchange={(e) => {
					pageSize = parseInt((e.target as HTMLSelectElement).value, 10) || 10;
					page = 1;
				}}
			>
				<option value={5}>5</option>
				<option value={10}>10</option>
				<option value={20}>20</option>
				<option value={50}>50</option>
			</select>
		</label>
		<div class="m-2 flex items-center gap-3 text-xs text-slate-300 sm:ml-auto">
			<button
				type="button"
				class="rounded-lg border border-white/15 bg-white/5 px-3 py-1.5 font-semibold text-white transition enabled:hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40"
				onclick={() => (page = Math.max(1, page - 1))}
				disabled={page === 1 || query.isFetching}
			>
				Prev
			</button>
			<span class="text-slate-200">
				Page {page} / {totalPages()}
			</span>
			<button
				type="button"
				class="rounded-lg border border-white/15 bg-white/5 px-3 py-1.5 font-semibold text-white transition enabled:hover:bg-white/10 disabled:cursor-not-allowed disabled:opacity-40"
				onclick={() => (page = Math.min(totalPages(), page + 1))}
				disabled={page >= totalPages() || query.isFetching}
			>
				Next
			</button>
		</div>
	</div>
</section>
