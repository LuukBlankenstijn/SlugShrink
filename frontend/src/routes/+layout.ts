import { browser } from '$app/environment';
import { QueryClient } from '@tanstack/svelte-query';
import type { LayoutLoad } from './$types';

export const prerender = false;

const browserQueryClient = browser ? new QueryClient() : null;

export const load: LayoutLoad = async () => {
	const queryClient = browserQueryClient ?? new QueryClient();
	return { queryClient };
};
