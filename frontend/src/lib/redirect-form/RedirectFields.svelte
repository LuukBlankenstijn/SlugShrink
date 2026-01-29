<script lang="ts">
	import QRCode from './QRCode.svelte';

	let {
		domainId = '',
		path = '',
		targetUrl = '',
		active = true,
		isViewMode = false,
		isExisting = false,
		isDomainsPending = false,
		domains = []
	}: {
		domainId?: string;
		path?: string;
		targetUrl?: string;
		active?: boolean;
		isViewMode?: boolean;
		isExisting?: boolean;
		isDomainsPending?: boolean;
		domains?: { id: string; domain: string }[];
	} = $props();

	const fieldBase =
		'w-full rounded-xl border px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20';

	const selectedDomain = $derived(
		domains.find((domain) => domain.id === domainId)?.domain ?? domainId
	);
</script>

{#if isExisting}
	<div class="flex flex-col items-start gap-1.5">
		<span class="text-sm text-slate-200">Status</span>
		{#if isViewMode}
			<div
				class={[
					'inline-flex items-center gap-2 rounded-full border px-3 py-1 text-xs font-semibold',
					active
						? 'border-emerald-400/30 bg-emerald-400/10 text-emerald-200'
						: 'border-slate-400/20 bg-white/5 text-slate-300'
				].join(' ')}
			>
				<span class={['h-2 w-2 rounded-full', active ? 'bg-emerald-300' : 'bg-slate-400'].join(' ')}
				></span>
				{active ? 'Active' : 'Paused'}
			</div>
		{:else}
			<label
				class="flex w-fit items-center gap-2 rounded-xl border border-white/10 bg-slate-950/30 px-3 py-2 text-sm text-slate-200 select-none"
			>
				<input type="checkbox" class="h-4 w-4 accent-cyan-400" bind:checked={active} />
				<span>Active</span>
			</label>
		{/if}
	</div>
{/if}

<label class="block space-y-1.5">
	<span class="text-sm text-slate-200">Domain</span>
	{#if isViewMode}
		<div
			class="rounded-lg bg-white/5 px-3 py-2 text-sm text-slate-200"
			title={selectedDomain || ''}
		>
			<span class="block truncate">{selectedDomain || '—'}</span>
		</div>
	{:else}
		<select
			class={[fieldBase, 'border-white/10 bg-slate-950/40'].join(' ')}
			bind:value={domainId}
			disabled={isDomainsPending || domains.length === 0}
		>
			{#if domains.length === 0}
				<option value="">No domains available</option>
			{:else}
				{#each domains as domain (domain.id)}
					<option value={domain.id}>{domain.domain}</option>
				{/each}
			{/if}
		</select>
	{/if}
</label>

<label class="block space-y-1.5">
	<span class="text-sm text-slate-200">Path</span>
	{#if isViewMode}
		<div class="rounded-lg bg-white/5 px-3 py-2 font-mono text-xs text-slate-200" title={path}>
			<span class="block truncate">{path || '—'}</span>
		</div>
	{:else}
		<input
			class={[fieldBase, 'border-white/10 bg-slate-950/40'].join(' ')}
			bind:value={path}
			placeholder="/promo"
			autocomplete="off"
		/>
	{/if}
</label>

<label class="block space-y-1.5">
	<span class="text-sm text-slate-200">Target URL</span>
	{#if isViewMode}
		<div class="rounded-lg bg-white/5 px-3 py-2 text-sm text-slate-200" title={targetUrl}>
			<span class="block truncate">{targetUrl || '—'}</span>
		</div>
	{:else}
		<input
			class={[fieldBase, 'border-white/10 bg-slate-950/40'].join(' ')}
			bind:value={targetUrl}
			placeholder="https://example.com/landing"
			autocomplete="off"
		/>
	{/if}
</label>

{#if isViewMode}
	<span class="text-sm text-slate-200">QR code</span>
	<div class="rounded-lg bg-white/5 px-3 py-2 text-sm text-slate-200" title="Click to copy">
		<QRCode domain={selectedDomain} {path} />
	</div>
{/if}
