import { makeApi } from '$lib/api.svelte';
import { queryKeys } from '$lib/queryKeys';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent, fetch }) => {
	const { queryClient } = await parent();
	const api = makeApi(fetch);

	await queryClient.prefetchQuery({
		queryKey: queryKeys.domains(),
		queryFn: () => api.getDomains({})
	});
};
