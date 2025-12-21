<script lang="ts">
	import { api } from '$lib/api.svelte';
	import { createQuery } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { MODAL_CONTEXT, type ModalControls } from '$lib/modal-context';
	import { queryKeys } from '$lib/queryKeys';
	import { authStatusQueryOptions } from '$lib/queries/authStatus';
	import { UserPermission } from '../../../gen/api/v1/auth_pb';

	const query = createQuery(() => ({
		queryKey: queryKeys.redirects(),
		queryFn: () => api.getRedirects({})
	}));

	const redirects = $derived(() => query.data?.data ?? []);
	const { openRedirectModal } = getContext<ModalControls>(MODAL_CONTEXT);

	const authStatus = createQuery(authStatusQueryOptions);
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
			<span class="font-semibold text-slate-100">{redirects().length}</span>
			{redirects().length === 1 ? 'redirect' : 'redirects'} found
		</p>
	</div>

	<div class="overflow-x-auto">
		<table class="min-w-full divide-y divide-white/10 text-sm">
			<thead class="bg-white/3">
				<tr class="text-left text-xs font-semibold tracking-wide text-slate-300 uppercase">
					<th class="px-4 py-3">Status</th>
					<th class="px-4 py-3">Domain</th>
					<th class="px-4 py-3">Path</th>
					<th class="px-4 py-3">Target</th>
					{#if authStatus.data && [UserPermission.PERMISSION_ADMIN, UserPermission.PERMISSION_SUPERUSER].includes(authStatus.data.permission)}
						<th class="px-4 py-3"></th>
					{/if}
				</tr>
			</thead>
			<tbody class="divide-y divide-white/10">
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
						{#if authStatus.data && [UserPermission.PERMISSION_ADMIN, UserPermission.PERMISSION_SUPERUSER].includes(authStatus.data.permission)}
							<td class="px-4 py-3 text-right">
								<button
									type="button"
									class="inline-flex items-center justify-center rounded-lg border border-white/10 bg-white/5 px-3 py-1.5 text-xs font-semibold text-white hover:bg-white/10 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
									onclick={() => openRedirectModal(r)}
								>
									Edit
								</button>
							</td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>
