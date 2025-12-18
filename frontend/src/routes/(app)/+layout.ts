import type { LayoutLoad } from './$types';

export const load: LayoutLoad = async ({ parent }) => {
	const { queryClient } = await parent();
	return { queryClient };
};
