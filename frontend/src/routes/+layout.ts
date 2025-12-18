import { browser } from '$app/environment';
import { QueryClient } from '@tanstack/svelte-query';

export const prerender = true;

export async function load() {
	const queryClient = new QueryClient({
		defaultOptions: {
			queries: {
				enabled: browser,
				staleTime: 60_000,
				refetchOnMount: false,
				refetchInterval: false
			}
		}
	});

	return { queryClient };
}
