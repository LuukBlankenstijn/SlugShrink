<script lang="ts">
	import './layout.css';
	import favicon from '$lib/assets/favicon.svg';
	import { QueryClientProvider } from '@tanstack/svelte-query';
	import type { LayoutData } from './$types';
	import type { Snippet } from 'svelte';
	import { onMount } from 'svelte';

	let { data, children }: { data: LayoutData; children: Snippet } = $props();

	onMount(() => {
		window.__TANSTACK_QUERY_CLIENT__ = data.queryClient;
	});
</script>

<svelte:head><link rel="icon" href={favicon} /></svelte:head>

<QueryClientProvider client={data.queryClient}>
	{@render children()}
</QueryClientProvider>
