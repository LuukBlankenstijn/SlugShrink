<script lang="ts">
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { untrack } from 'svelte';
	import { authStatusQueryKey } from '$lib/queries/authStatus';
	import { queryKeys } from '$lib/queryKeys';
	import { authlessConfigMutationOptions } from '$lib/mutationOptions';

	const { saveTrigger } = $props();
	const queryclient = useQueryClient();
	let lastProcessedTrigger = $state(0);

	const mutation = createMutation(() => ({
		...authlessConfigMutationOptions(),
		onSuccess: () => {
			queryclient.invalidateQueries({ queryKey: queryKeys.authConfig() });
			queryclient.invalidateQueries({ queryKey: authStatusQueryKey });
		},
	}));

	$effect(() => {
		if (saveTrigger <= lastProcessedTrigger || mutation.isPending) return;
		lastProcessedTrigger = saveTrigger;
		untrack(() => {
			mutation.mutate();
		});
	});
</script>

<div
	class="flex h-full items-center justify-center rounded-xl border border-dashed border-white/10 bg-slate-950/30 px-4 py-8 text-sm text-slate-300"
>
	No settings required for this method.
</div>
