<script lang="ts">
	type Redirect = {
		id: string;
		domain_id: string;
		path: string;
		target_url: string;
		active: boolean;
	};

	const domainsById: Record<string, string> = {
		'2bdb5a11-61d3-4f12-9bb8-8d20a4c7d663': 'gewi.sh',
		'1e64b1ea-0af9-4a75-9b79-179b2a0e6e7a': 'app.gewi.sh',
		'e2f72e31-8e6d-4a3f-9c32-5e8c8d1b1f0a': 'api.gewi.sh'
	};

	const redirects: Redirect[] = [
		{
			id: '6e5d2e3a-4c92-4b2a-9e56-6d5c902efc9a',
			domain_id: '2bdb5a11-61d3-4f12-9bb8-8d20a4c7d663',
			path: '/pricing',
			target_url: 'https://gewi.sh/pricing',
			active: true
		},
		{
			id: '0f0a8b65-2e90-42b9-8a45-d5c90c9b6c0e',
			domain_id: '2bdb5a11-61d3-4f12-9bb8-8d20a4c7d663',
			path: '/docs',
			target_url: 'https://docs.gewi.sh/',
			active: true
		},
		{
			id: 'e9c2d1c1-6c22-43b8-8f97-1e03bd7b8d69',
			domain_id: '2bdb5a11-61d3-4f12-9bb8-8d20a4c7d663',
			path: '/discord',
			target_url: 'https://discord.gg/example',
			active: true
		},
		{
			id: '7b72f6a4-2b69-4a25-b4f3-0c8b23b5b2db',
			domain_id: '2bdb5a11-61d3-4f12-9bb8-8d20a4c7d663',
			path: '/gh',
			target_url: 'https://github.com/gewish/gewish',
			active: true
		},
		{
			id: '93a3ef6c-95b7-42a1-9f4d-6edb47c8ce2c',
			domain_id: '1e64b1ea-0af9-4a75-9b79-179b2a0e6e7a',
			path: '/support',
			target_url: 'https://gewi.sh/support',
			active: true
		},
		{
			id: '60e1e7d2-8889-4c2e-9e5c-fd9f0fd4a0aa',
			domain_id: '1e64b1ea-0af9-4a75-9b79-179b2a0e6e7a',
			path: '/status',
			target_url: 'https://status.gewi.sh/',
			active: false
		},
		{
			id: '6d73bff8-1df2-4a6d-a7c2-1b9b7f3a4f84',
			domain_id: '1e64b1ea-0af9-4a75-9b79-179b2a0e6e7a',
			path: '/blog',
			target_url: 'https://gewi.sh/blog',
			active: true
		},
		{
			id: '5b6d7c88-0e7b-4f8f-90c7-3b4c5d6e7f81',
			domain_id: 'e2f72e31-8e6d-4a3f-9c32-5e8c8d1b1f0a',
			path: '/waitlist',
			target_url: 'https://gewi.sh/waitlist',
			active: true
		},
		{
			id: 'f3a9c2a0-1f6e-4c2a-9ef8-5f1b0f7b2a7c',
			domain_id: 'e2f72e31-8e6d-4a3f-9c32-5e8c8d1b1f0a',
			path: '/v1',
			target_url: 'https://api.gewi.sh/v1',
			active: false
		},
		{
			id: 'a1c5c22a-0f62-4e1d-9b25-25fcb18f2b12',
			domain_id: 'e2f72e31-8e6d-4a3f-9c32-5e8c8d1b1f0a',
			path: '/terms',
			target_url: 'https://gewi.sh/terms',
			active: true
		}
	];

	let query = $state('');
	let includeInactive = $state(true);

	const filtered = $derived.by(() => {
		const q = query.trim().toLowerCase();

		return redirects.filter((r) => {
			if (!includeInactive && !r.active) return false;
			if (!q) return true;
			const domain = (domainsById[r.domain_id] ?? '').toLowerCase();
			return (
				r.path.toLowerCase().includes(q) ||
				r.target_url.toLowerCase().includes(q) ||
				r.id.toLowerCase().includes(q) ||
				r.domain_id.toLowerCase().includes(q) ||
				domain.includes(q)
			);
		});
	});
</script>

<div class="min-h-svh bg-slate-950 text-slate-100">
	<div class="mx-auto flex max-w-6xl flex-col gap-6 px-4 py-10">
		<header class="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
			<div class="space-y-1">
				<h1 class="text-2xl font-semibold tracking-tight">Redirects</h1>
				<div
					class="h-px w-24 bg-linear-to-r from-cyan-400/70 via-blue-500/70 to-violet-500/70"
				></div>
			</div>

			<div class="flex items-center gap-3">
				<button
					type="button"
					class="inline-flex items-center justify-center rounded-xl border border-white/10 bg-white/5 px-4 py-2 text-sm font-semibold text-white hover:bg-white/10 focus:outline-none focus-visible:ring-2 focus-visible:ring-cyan-400/40"
				>
					New
				</button>
			</div>
		</header>

		<section
			class="overflow-hidden rounded-2xl border border-white/10 bg-white/3"
			aria-label="Redirects table"
		>
			<div
				class="flex flex-col gap-3 border-b border-white/10 p-4 sm:flex-row sm:items-center sm:justify-between"
			>
				<div class="flex flex-col gap-2 sm:flex-row sm:items-center">
					<label class="relative w-full sm:w-88">
						<span class="sr-only">Search</span>
						<input
							class="w-full rounded-xl border border-white/10 bg-slate-950/40 px-3 py-2 text-sm text-slate-100 ring-0 outline-none placeholder:text-slate-400 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/15 focus:outline-none"
							placeholder="Search path, target URL, id…"
							bind:value={query}
						/>
					</label>

					<label
						class="inline-flex items-center gap-2 rounded-xl border border-white/10 bg-slate-950/30 px-3 py-2 text-sm text-slate-200 select-none"
					>
						<input type="checkbox" class="h-4 w-4 accent-cyan-400" bind:checked={includeInactive} />
						<span>Include inactive</span>
					</label>
				</div>

				<p class="text-sm text-slate-300">
					<span class="font-semibold text-slate-100">{filtered.length}</span>
					{filtered.length === 1 ? 'redirect' : 'redirects'}
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
						</tr>
					</thead>
					<tbody class="divide-y divide-white/10">
						{#each filtered as r (r.id)}
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
										<span
											class="truncate text-slate-200"
											title={domainsById[r.domain_id] ?? r.domain_id}
										>
											{domainsById[r.domain_id] ?? r.domain_id}
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
									<span class="block max-w-136 truncate text-slate-200" title={r.target_url}>
										{r.target_url}
									</span>
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</section>
	</div>
</div>
