import { api } from '$lib/api.svelte';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent }) => {
	const { queryClient } = await parent();

	await queryClient.prefetchQuery({
		queryKey: ['domains'],
		queryFn: () => api.getDomains({})
	});
};
