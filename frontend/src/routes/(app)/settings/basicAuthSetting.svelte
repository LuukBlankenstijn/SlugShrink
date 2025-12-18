<script lang="ts">
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { api } from '$lib/api.svelte';
	import { createMutation, useQueryClient } from '@tanstack/svelte-query';
	import { untrack } from 'svelte';
	import { authStatusQueryKey } from '$lib/queries/authStatus';
	import { queryKeys } from '$lib/queryKeys';

	const { saveTrigger } = $props();
	const queryclient = useQueryClient();
	let lastProcessedTrigger = saveTrigger;

	let password = $state('');
	let repeatedPassword = $state('');
	let validationError = $state<string | null>(null);
	let serverError = $state<string | null>(null);
	const MIN_PASSWORD_LENGTH = 8;
	let showPassword = $state(false);
	let showRepeatedPassword = $state(false);

	const errorMessage = $derived(validationError ?? serverError);

	const mutation = createMutation(() => ({
		mutationFn: () =>
			api.setAuthConfig({
				config: {
					case: 'basicAuth',
					value: {
						password: password,
						repeatedPassword: repeatedPassword
					}
				}
			}),
		onError: (err: unknown) => {
			serverError = err instanceof Error ? err.message : 'Something went wrong';
		},
		onSuccess: () => {
			validationError = null;
			serverError = null;
			queryclient.invalidateQueries({ queryKey: queryKeys.authConfig() });
			queryclient.invalidateQueries({ queryKey: authStatusQueryKey });
			goto(resolve('/login'));
		},
		retry: false
	}));

	function validatePasswords() {
		validationError = null;
		serverError = null;
		const hasUpper = /[A-Z]/.test(password);
		const hasLower = /[a-z]/.test(password);
		const hasNumber = /\d/.test(password);
		const hasSpecial = /[^A-Za-z0-9]/.test(password);

		if (!password || !repeatedPassword) {
			validationError = 'Enter your new password in both fields.';
			return false;
		}

		if (password.length < MIN_PASSWORD_LENGTH) {
			validationError = `Password must be at least ${MIN_PASSWORD_LENGTH} characters long.`;
			return false;
		}

		if (!(hasUpper && hasLower && hasNumber && hasSpecial)) {
			validationError =
				'Password needs at least one uppercase letter, one lowercase letter, one number, and one special character.';
			return false;
		}

		if (password !== repeatedPassword) {
			validationError = 'Passwords do not match.';
			return false;
		}

		return true;
	}

	$effect(() => {
		if (saveTrigger > lastProcessedTrigger && !mutation.isPending) {
			lastProcessedTrigger = saveTrigger;
			untrack(() => {
				if (validatePasswords()) {
					mutation.mutate();
				}
			});
		}
	});

	$effect(() => {
		// Clear errors as the user edits either field.
		validationError = null;
		serverError = null;
	});
</script>

<div class="space-y-4 rounded-xl border border-white/10 bg-slate-950/40 p-4">
	<h4 class="whitespace-nowrap">Update Password:</h4>
	<div class="space-y-1">
		<label for="password-input" class="text-sm text-slate-200">Password</label>
		<div class="relative">
			<input
				bind:value={password}
				id="password-input"
				type={showPassword ? 'text' : 'password'}
				class="w-full rounded-lg border border-white/10 bg-slate-900/70 px-3 py-2 pr-12 text-sm text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
				placeholder="Enter a strong password"
			/>
			<button
				type="button"
				class="absolute inset-y-0 right-2 my-1 rounded-lg px-3 text-xs font-semibold text-slate-200 hover:bg-white/10 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
				onclick={() => (showPassword = !showPassword)}
				aria-label={showPassword ? 'Hide password' : 'Show password'}
			>
				{showPassword ? 'Hide' : 'Show'}
			</button>
		</div>
	</div>
	<div class="space-y-1">
		<label for="password-repeat" class="text-sm text-slate-200">Repeat password</label>
		<div class="relative">
			<input
				bind:value={repeatedPassword}
				id="password-repeat"
				type={showRepeatedPassword ? 'text' : 'password'}
				class="w-full rounded-lg border border-white/10 bg-slate-900/70 px-3 py-2 pr-12 text-sm text-slate-50 outline-none placeholder:text-slate-500 focus:border-cyan-400/40 focus:ring-2 focus:ring-cyan-400/20"
				placeholder="Repeat password"
			/>
			<button
				type="button"
				class="absolute inset-y-0 right-2 my-1 rounded-lg px-3 text-xs font-semibold text-slate-200 hover:bg-white/10 focus-visible:ring-2 focus-visible:ring-cyan-400/40 focus-visible:outline-none"
				onclick={() => (showRepeatedPassword = !showRepeatedPassword)}
				aria-label={showRepeatedPassword ? 'Hide repeated password' : 'Show repeated password'}
			>
				{showRepeatedPassword ? 'Hide' : 'Show'}
			</button>
		</div>
	</div>

	{#if errorMessage}
		<p class="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">
			{errorMessage}
		</p>
	{/if}

	<p class="text-xs text-slate-400">
		Use at least 8 characters with uppercase, lowercase, number, and special symbols.
	</p>
</div>
