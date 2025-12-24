<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { authStatusQueryOptions } from '$lib/queries/authStatus';
	import { queryKeys } from '$lib/queryKeys';
	import { loginMutationOptions } from '$lib/mutationOptions';
	const queryClient = useQueryClient();

	let password = $state('');
	let error = $state('');

	const authStatus = createQuery(authStatusQueryOptions);

	const login = createMutation(() => ({
		...loginMutationOptions(),
		onSuccess: () => {
			queryClient.setQueryData(queryKeys.authStatus(), () => true);
			error = '';
		},
		onError: (err: unknown) => {
			error = err instanceof Error ? err.message : '';
		}
	}));

	const pending = $derived(login.isPending);

	$effect(() => {
		if (authStatus.data) {
			goto(resolve('/(app)/redirects'));
		}
	});

	function handleSubmit(event: Event) {
		event.preventDefault();
		error = '';
		login.mutate({ password });
	}
</script>

<svelte:head>
	<title>Sign in · slugshrink</title>
</svelte:head>

<div class="min-h-svh bg-slate-950 text-slate-100">
	<div class="mx-auto flex max-w-5xl flex-col gap-10 px-4 py-12">
		<div class="space-y-3 text-center">
			<p class="inline-flex items-center justify-center gap-2 text-sm text-slate-300">
				<span
					class="h-2 w-2 rounded-full bg-linear-to-r from-cyan-400 via-blue-500 to-violet-500 shadow-[0_0_18px_rgba(34,211,238,0.45)]"
				></span>
				<span
					class="bg-linear-to-r from-cyan-400 via-blue-500 to-violet-500 bg-clip-text font-semibold text-transparent"
				>
					slugshrink
				</span>
			</p>
			<h1 class="text-3xl font-semibold tracking-tight">Welcome back</h1>
		</div>

		<div
			class="mx-auto w-full max-w-xl overflow-hidden rounded-2xl border border-white/10 bg-white/5 shadow-[0_18px_80px_rgba(0,0,0,0.35)]"
		>
			<div class="border-b border-white/10 px-6 py-4">
				<p class="text-sm font-semibold text-slate-100">Sign in</p>
			</div>

			<form class="space-y-5 px-6 py-6" onsubmit={handleSubmit}>
				<label class="block space-y-1">
					<span class="text-sm text-slate-200">Password</span>
					<input
						bind:value={password}
						class="w-full rounded-xl border border-white/10 bg-slate-950/50 px-3 py-2 text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20 disabled:opacity-60"
						type="password"
						placeholder="••••••••"
						autocomplete="current-password"
						required
						disabled={pending}
					/>
				</label>

				{#if error}
					<p class="text-sm text-rose-300">{error}</p>
				{/if}

				<button
					type="submit"
					class="mt-2 w-full rounded-xl bg-linear-to-r from-cyan-500 via-blue-600 to-violet-600 px-4 py-2 text-sm font-semibold text-white shadow-[0_12px_32px_rgba(59,130,246,0.25)] transition enabled:translate-y-px enabled:hover:brightness-110 enabled:focus-visible:ring-2 enabled:focus-visible:ring-cyan-400/50 enabled:focus-visible:outline-none disabled:cursor-not-allowed disabled:opacity-80"
					disabled={pending}
				>
					{pending ? 'Signing in…' : 'Sign in'}
				</button>
			</form>
		</div>
	</div>
</div>
