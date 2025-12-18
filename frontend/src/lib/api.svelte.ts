import { createClient } from '@connectrpc/connect';
import { ApiService } from '../gen/api/v1/api_pb';
import { createConnectTransport } from '@connectrpc/connect-web';

const transport = createConnectTransport({
	baseUrl: '/',
	fetch
});

export const api = createClient(ApiService, transport);
