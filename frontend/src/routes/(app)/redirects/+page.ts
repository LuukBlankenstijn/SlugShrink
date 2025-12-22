import { makeApi } from '$lib/api.svelte';
import { redirectsQueryOptions } from '$lib/queryOptions';
import { browser } from '$app/environment';
import type { PageLoad } from './$types';

export const load: PageLoad = async ({ parent, fetch }) => {
  const { queryClient } = await parent();
  const api = makeApi(fetch);

  if (!browser) return;
  await queryClient.prefetchQuery(redirectsQueryOptions({ page: 1, pageSize: 10, search: "" }, api));
};
