<script lang="ts">
	import { api } from '$lib/api.svelte';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { untrack } from 'svelte';
	import { authStatusQueryKey } from '$lib/queries/authStatus';
	import { queryKeys } from '$lib/queryKeys';

	const { saveTrigger } = $props();
	const queryclient = useQueryClient();
	let lastProcessedTrigger = saveTrigger;

	const mutation = createMutation(() => ({
		mutationFn: () =>
			api.setAuthConfig({
				config: {
					case: 'authless',
					value: {}
				}
			}),
		onSuccess: () => {
			queryclient.invalidateQueries({ queryKey: queryKeys.authConfig() });
			queryclient.invalidateQueries({ queryKey: authStatusQueryKey });
		},
		retry: false
	}));

	$effect(() => {
		if (saveTrigger > lastProcessedTrigger && !mutation.isPending) {
			lastProcessedTrigger = saveTrigger;
			untrack(() => {
				mutation.mutate();
			});
		}
	});
</script>

<div
	class="flex h-full items-center justify-center rounded-xl border border-dashed border-white/10 bg-slate-950/30 px-4 py-8 text-sm text-slate-300"
>
	No settings required for this method.
</div>
