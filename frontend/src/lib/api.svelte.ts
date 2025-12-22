import { createClient } from '@connectrpc/connect';
import { ApiService } from '../gen/api/v1/api_pb';
import { createConnectTransport } from '@connectrpc/connect-web';

export const makeApi = (fetchFn: typeof fetch) => {
	const withCredentials = (input: RequestInfo | URL, init?: RequestInit) =>
		fetchFn(input, { credentials: 'include', ...init });

	const transport = createConnectTransport({
		baseUrl: '/',
		fetch: withCredentials
	});

	return createClient(ApiService, transport);
};

export type ApiClient = ReturnType<typeof makeApi>;

export const api = makeApi(fetch);
