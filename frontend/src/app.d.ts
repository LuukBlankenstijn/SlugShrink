declare global {
	namespace App {
		// interface Error {}
		// interface Locals {}
		// interface PageData {}
		// interface PageState {}
		// interface Platform {}
	}

	interface Window {
		__TANSTACK_QUERY_CLIENT__?: import('@tanstack/query-core').QueryClient;
	}
}

export {};
