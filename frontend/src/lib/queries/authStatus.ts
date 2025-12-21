import { api } from '$lib/api.svelte';
import { queryKeys } from '$lib/queryKeys';

export const authStatusQueryKey = queryKeys.authStatus();

export const authStatusQueryOptions = () => ({
  queryKey: authStatusQueryKey,
  queryFn: async () => {
    const result = await api.getAuthStatus({});
    if (result.status.case === "authenticated") {
      return result.status.value
    }
    return false
  },
  retry: false
});
