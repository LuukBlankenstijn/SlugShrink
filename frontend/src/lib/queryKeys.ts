export const queryKeys = {
  allRedirects: () => ['redirects'],
  redirects: (page: number, pageSize: number) => ['redirects', page, pageSize] as const,
  domains: () => ['domains'] as const,
  authConfig: () => ['authConfig'] as const,
  authStatus: () => ['authStatus'] as const
};
