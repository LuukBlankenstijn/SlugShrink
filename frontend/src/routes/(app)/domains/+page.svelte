<script lang="ts">
	import { createQuery } from '@tanstack/svelte-query';
	import { getContext } from 'svelte';
	import { MODAL_CONTEXT, type ModalControls } from '$lib/modal-context';
	import { authStatusQueryOptions } from '$lib/queries/authStatus';
	import { resolve } from '$app/paths';
	import { goto } from '$app/navigation';
	import { UserPermission } from '../../../gen/api/v1/auth_pb';
	import { hasPermission } from '$lib/auth';
	import { domainsQueryOptions } from '$lib/queryOptions';

	const authStatus = createQuery(authStatusQueryOptions);

	$effect(() => {
		if (
			!hasPermission(
				authStatus.data,
				UserPermission.PERMISSION_ADMIN,
				UserPermission.PERMISSION_SUPERUSER
			)
		) {
			goto(resolve('/redirects'));
		}
	});

	const query = createQuery(() => domainsQueryOptions());

	const domains = $derived(() => query.data?.data ?? []);
	const { openDomainModal } = getContext<ModalControls>(MODAL_CONTEXT);
</script>

<svelte:head>
	<title>Domains · gewi.sh</title>
</svelte:head>

<section
	class="overflow-hidden rounded-2xl border border-white/10 bg-white/3"
	aria-label="Domains table"
>
	<div
		class="flex flex-col gap-3 border-b border-white/10 p-4 sm:flex-row sm:items-center sm:justify-between"
	>
		<p class="text-sm text-slate-300">
			<span class="font-semibold text-slate-100">{domains().length}</span>
			{domains().length === 1 ? 'domain' : 'domains'} found
		</p>
	</div>

	<div class="overflow-x-auto">
		<table class="min-w-full divide-y divide-white/10 text-sm">
			<thead class="bg-white/3">
				<tr class="text-left text-xs font-semibold tracking-wide text-slate-300 uppercase">
					<th class="px-4 py-3">Name</th>
					<th class="px-4 py-3">Domain</th>
					<th class="px-4 py-3"></th>
				</tr>
			</thead>
			<tbody class="divide-y divide-white/10">
				{#each domains() as d (d.id)}
					<tr class="group hover:bg-white/4">
						<td class="px-4 py-3">
							<span class="inline-flex max-w-[16rem] items-center gap-2">
								<span class="truncate text-slate-200" title={d.name}>
									{d.name}
								</span>
							</span>
						</td>
						<td class="px-4 py-3">
							<span class="inline-flex max-w-[16rem] items-center gap-2">
								<span class="truncate text-slate-200" title={d.domain}>
									{d.domain}
								</span>
							</span>
						</td>
						<td class="px-4 py-3 text-right">
							<button
								type="button"
								class="inline-flex items-center justify-center rounded-lg border border-white/15 bg-white/8 px-3 py-1.5 text-xs font-semibold text-white transition hover:border-white/25 hover:bg-white/12 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
								onclick={() => openDomainModal(d)}
							>
								Edit
							</button>
						</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
</section>
