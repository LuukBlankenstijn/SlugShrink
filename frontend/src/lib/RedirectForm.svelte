<script lang="ts">
	import { createMutation, createQuery, useQueryClient } from '@tanstack/svelte-query';
	import { api } from './api.svelte';
	import RedirectFields from '$lib/redirect-form/RedirectFields.svelte';
	import RedirectFormActions from '$lib/redirect-form/RedirectFormActions.svelte';
	import { queryKeys } from '$lib/queryKeys';
	import { authStatusQueryOptions } from '$lib/queries/authStatus';
	import type { FullRedirect } from '../gen/api/v1/redirect_pb';
	import { domainsQueryOptions } from '$lib/queryOptions';
	import { deleteRedirectMutationOptions, saveRedirectMutationOptions } from '$lib/mutationOptions';
	import { UserPermission } from '../gen/api/v1/auth_pb';

	interface Props {
		readonly onClose: () => void;
		readonly current?: FullRedirect;
	}

	const { onClose, current }: Props = $props();

	const qc = useQueryClient();
	const domainsQuery = createQuery(() => domainsQueryOptions(api));
	const authStatus = createQuery(authStatusQueryOptions);

	const domains = $derived(() => domainsQuery.data?.data ?? []);

	let domainId = $state('');
	let path = $state('');
	let targetUrl = $state('');
	let active = $state(true);
	let error = $state<string | null>(null);

	const isExisting = $derived(() => Boolean(current));
	// svelte-ignore state_referenced_locally
	let isViewMode = $state(Boolean(current));
	const canEdit = $derived(() => {
		if (!current || !authStatus.data) {
			return false;
		}
		if (
			authStatus.data.permission === UserPermission.PERMISSION_ADMIN ||
			authStatus.data.permission === UserPermission.PERMISSION_SUPERUSER
		) {
			return true;
		}
		return authStatus.data.userId !== undefined && authStatus.data.userId === current.creator;
	});

	const close = () => {
		onClose();
		domainId = '';
		path = '';
		targetUrl = '';
		active = true;
		error = null;
	};

	const cancelEdit = () => {
		if (!current) {
			return;
		}
		domainId = current.domainId;
		path = current.path;
		targetUrl = current.targetUrl;
		active = current.active;
		error = null;
		isViewMode = true;
	};

	$effect(() => {
		if (current) {
			domainId = current.domainId;
			path = current.path;
			targetUrl = current.targetUrl;
			active = current.active;
			isViewMode = true;
		} else {
			domainId = '';
			path = '';
			targetUrl = '';
			active = true;
			isViewMode = false;
		}
		error = null;
	});

	$effect(() => {
		if (!domainId && domains().length > 0) {
			domainId = domains()[0].id;
		}
	});

	const saveRedirect = createMutation(() => ({
		...saveRedirectMutationOptions(api, current),
		onSuccess: async () => {
			qc.invalidateQueries({ queryKey: queryKeys.allRedirects() });
			close();
		},
		onError: (e) => {
			error = e instanceof Error ? e.message : 'Something went wrong';
		}
	}));

	const deleteRedirect = createMutation(() => ({
		...deleteRedirectMutationOptions(api),
		onSuccess: async () => {
			qc.invalidateQueries({ queryKey: queryKeys.allRedirects() });
			close();
		},
		onError: (e) => {
			error = e instanceof Error ? e.message : 'Something went wrong';
		}
	}));
</script>

<form
	class="space-y-5"
	onsubmit={(e) => {
		e.preventDefault();
		if (isExisting() && isViewMode) {
			if (!canEdit()) {
				error = 'You do not have permission to edit this redirect.';
				return;
			}
			isViewMode = false;
			return;
		}
		saveRedirect.mutate({ domainId, path, targetUrl, active });
	}}
>
	<RedirectFields
		bind:domainId
		bind:path
		bind:targetUrl
		bind:active
		{isViewMode}
		isExisting={isExisting()}
		isDomainsPending={domainsQuery.isPending}
		domains={domains()}
	/>

	{#if error}
		<p class="rounded-lg border border-red-500/30 bg-red-500/10 px-3 py-2 text-sm text-red-100">
			{error}
		</p>
	{/if}

	<RedirectFormActions
		{current}
		isExisting={isExisting()}
		{isViewMode}
		canEdit={canEdit()}
		isSaving={saveRedirect.isPending}
		canCreate={domains().length > 0}
		isDeleting={deleteRedirect.isPending}
		onCancel={cancelEdit}
		onDelete={(id) => deleteRedirect.mutate({ id })}
	/>
</form>
