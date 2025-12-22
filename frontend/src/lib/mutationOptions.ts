import { api, type ApiClient } from '$lib/api.svelte';
import type { Domain } from '../gen/api/v1/domain_pb';
import type { FullRedirect } from '../gen/api/v1/redirect_pb';

export const saveDomainMutationOptions = (client: ApiClient = api, current?: Domain) => ({
	mutationFn: ({ name, domain }: { name: string; domain: string }) =>
		current ? client.putDomain({ id: current.id, name, domain }) : client.createDomain({ name, domain })
});

export const deleteDomainMutationOptions = (client: ApiClient = api) => ({
	mutationFn: ({ id }: { id: string }) => client.deleteDomain({ id })
});

export const saveRedirectMutationOptions = (client: ApiClient = api, current?: FullRedirect) => ({
	mutationFn: ({
		domainId,
		path,
		targetUrl,
		active
	}: {
		domainId: string;
		path: string;
		targetUrl: string;
		active: boolean;
	}) =>
		current
			? client.putRedirect({ id: current.id, domainId, path, targetUrl, active })
			: client.createRedirect({ domainId, path, targetUrl, active })
});

export const deleteRedirectMutationOptions = (client: ApiClient = api) => ({
	mutationFn: ({ id }: { id: string }) => client.deleteRedirect({ id })
});

export const basicAuthConfigMutationOptions = (
	client: ApiClient = api,
	{ password, repeatedPassword }: { password: string; repeatedPassword: string }
) => ({
	mutationFn: () =>
		client.setAuthConfig({
			config: {
				case: 'basicAuth',
				value: {
					password,
					repeatedPassword
				}
			}
		}),
	retry: false
});

export const authlessConfigMutationOptions = (client: ApiClient = api) => ({
	mutationFn: () =>
		client.setAuthConfig({
			config: {
				case: 'authless',
				value: {}
			}
		}),
	retry: false
});

export const proxyAuthConfigMutationOptions = (
	client: ApiClient = api,
	{
		groupHeader,
		userIdHeader,
		superUserGroups,
		adminGroups,
		groupSeparator
	}: {
		groupHeader: string;
		userIdHeader: string;
		superUserGroups: string[];
		adminGroups: string[];
		groupSeparator: string;
	}
) => ({
	mutationFn: () =>
		client.setAuthConfig({
			config: {
				case: 'proxyAuth',
				value: {
					groupHeader,
					userIdHeader,
					superUserGroups,
					adminGroups,
					groupsSeparator: groupSeparator ? groupSeparator : undefined
				}
			}
		}),
	retry: false
});

export const loginMutationOptions = (client: ApiClient = api) => ({
	mutationFn: ({ password }: { password: string }) => client.login({ password })
});
