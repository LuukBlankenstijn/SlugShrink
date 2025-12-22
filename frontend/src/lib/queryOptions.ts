import { api, type ApiClient } from '$lib/api.svelte';
import { browser } from '$app/environment';
import { queryKeys } from '$lib/queryKeys';
import type { FullRedirect } from '../gen/api/v1/redirect_pb';

export const domainsQueryOptions = (client: ApiClient = api) => ({
  queryKey: queryKeys.domains(),
  queryFn: () => client.getDomains({}),
  enabled: browser
});

export const redirectsQueryOptions = (
  {
    page,
    pageSize,
    search
  }: {
    page: number;
    pageSize: number;
    search: string;
  },
  client: ApiClient = api
) => ({
  queryKey: queryKeys.redirects(page, pageSize, search),
  queryFn: () =>
    client.getRedirects({
      page,
      pagesize: pageSize,
      search
    }),
  enabled: browser,
  keepPreviousData: true,
  placeholderData: (prev: { data: FullRedirect[]; total: number, $typeName: 'api.v1.RedirectsResponse' } | undefined) =>
    prev ?? { $typeName: 'api.v1.RedirectsResponse' as const, data: [], total: 0 }
});

export const authConfigQueryOptions = (client: ApiClient = api) => ({
  queryKey: queryKeys.authConfig(),
  queryFn: () => client.getAuthConfig({}),
  enabled: browser
});
