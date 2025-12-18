import { api } from '$lib/api.svelte';
import { queryKeys } from '$lib/queryKeys';

export const authStatusQueryKey = queryKeys.authStatus();

export const authStatusQueryOptions = () => ({
	queryKey: authStatusQueryKey,
	queryFn: async () => {
		const result = await api.getAuthStatus({});

		if (typeof result === 'object' && result !== null && 'authenticated' in result) {
			// Backward compatibility for the previous shape.
			return Boolean((result as { authenticated?: boolean }).authenticated);
		}

		return Boolean(result);
	},
	retry: false
});
