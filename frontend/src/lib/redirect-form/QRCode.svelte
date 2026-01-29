<script lang="ts">
	import { createQrPngDataUrl } from '@svelte-put/qr';
	import { qr } from '@svelte-put/qr/svg';

	let {
		domain,
		path
	}: {
		domain: string;
		path: string;
	} = $props();

	const fullUrl = $derived.by(() => {
		const protocol = typeof window !== 'undefined' ? window.location.protocol : 'https:';

		return `${protocol}//${domain}${path}`;
	});

	const copyPng = async () => {
		try {
			const png = await createQrPngDataUrl({
				data: fullUrl,
				width: 1000,
				height: 1000,
				backgroundFill: '#FFFFFF',
				margin: 3,
				moduleFill: '#000',
				anchorOuterFill: '#000',
				anchorInnerFill: '#000'
			});

			const res = await fetch(png);
			const blob = await res.blob();

			const item = new ClipboardItem({ 'image/png': blob });

			await navigator.clipboard.write([item]);
		} catch (err) {
			console.error('Failed to copy qr code: ', err);
		}
	};
</script>

<button
	type="button"
	onclick={copyPng}
	aria-label="Copy QR code to clipboard"
	class="group relative block rounded-lg border border-white/10 bg-white p-2
           transition-all duration-150 ease-in-out
           hover:scale-[1.02]
           focus-visible:ring-2
           focus-visible:ring-cyan-400 focus-visible:outline-none active:scale-100"
>
	<svg
		class="w-40"
		use:qr={{
			data: fullUrl,
			margin: 3,
			moduleFill: '#000',
			anchorOuterFill: '#000',
			anchorInnerFill: '#000'
		}}
	></svg>
</button>
