export const queryKeys = {
  allRedirects: () => ['redirects'],
  redirects: (page: number, pageSize: number, search: string) => ['redirects', page, pageSize, search] as const,
  domains: () => ['domains'] as const,
  authConfig: () => ['authConfig'] as const,
  authStatus: () => ['authStatus'] as const
};
