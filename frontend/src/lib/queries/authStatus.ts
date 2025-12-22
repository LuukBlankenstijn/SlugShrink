import { api, type ApiClient } from '$lib/api.svelte';
import { browser } from '$app/environment';
import { queryKeys } from '$lib/queryKeys';

export const authStatusQueryKey = queryKeys.authStatus();

export const authStatusQueryOptions = (client: ApiClient = api) => ({
	queryKey: authStatusQueryKey,
	queryFn: async () => {
		const result = await client.getAuthStatus({});
		if (result.status.case === 'authenticated') {
			return result.status.value;
		}
		return false;
	},
	enabled: browser,
	retry: false
});
