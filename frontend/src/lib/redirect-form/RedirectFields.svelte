<script lang="ts">
	import QRCode from './QRCode.svelte';

	let {
		domainId = $bindable(''),
		path = $bindable(''),
		targetUrl = $bindable(''),
		active = $bindable(true),
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

	let copied = $state(false);
	let copyTimeout: ReturnType<typeof setTimeout>;

	function copyTargetUrl() {
		navigator.clipboard.writeText(targetUrl);
		copied = true;
		clearTimeout(copyTimeout);
		copyTimeout = setTimeout(() => (copied = false), 2000);
	}
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
		<div class="flex items-center gap-2 rounded-lg bg-white/5 px-3 py-2 text-sm text-slate-200" title={targetUrl}>
			<span class="block flex-1 truncate">{targetUrl || '—'}</span>
			{#if targetUrl}
				<button
					type="button"
					class="shrink-0 cursor-pointer text-slate-400 transition-colors hover:text-white"
					title={copied ? 'Copied!' : 'Copy target URL'}
					onclick={copyTargetUrl}
				>
					{#if copied}
						<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-emerald-400" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<polyline points="20 6 9 17 4 12"></polyline>
						</svg>
					{:else}
						<svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
							<rect x="9" y="9" width="13" height="13" rx="2" ry="2"></rect>
							<path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
						</svg>
					{/if}
				</button>
			{/if}
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
