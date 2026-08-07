import {
	createRouter as createTanStackRouter,
	useRouterState,
} from "@tanstack/react-router";
import { setupRouterSsrQueryIntegration } from "@tanstack/react-router-ssr-query";
import { RoutingPending } from "./components/shared";
import { getContext } from "./integrations/tanstack-query/root-provider";
import { routeTree } from "./routeTree.gen";

// Names the path that failed to match — the default is a bare <p>Not Found</p>,
// which hides the one fact you need to debug it.
function NotFound() {
	const pathname = useRouterState({ select: (s) => s.location.pathname });
	return (
		<div className="p-8">
			<h1 className="text-lg font-semibold">Not found</h1>
			<p className="text-muted-foreground mt-1 text-sm">
				No route matched <code>{pathname}</code>
			</p>
		</div>
	);
}

export function getRouter() {
	const context = getContext();

	const router = createTanStackRouter({
		routeTree,
		context,
		scrollRestoration: true,
		defaultPreload: "intent",
		defaultPreloadStaleTime: 0,
		defaultNotFoundComponent: NotFound,
		defaultPendingComponent: RoutingPending,
		// Below this, a cache-hit (e.g. post-login) resolves before it'd ever show.
		defaultPendingMs: 150,
		// Once shown, hold it this long so it doesn't flicker in and out.
		defaultPendingMinMs: 300,
	});

	setupRouterSsrQueryIntegration({ router, queryClient: context.queryClient });

	return router;
}

declare module "@tanstack/react-router" {
	interface Register {
		router: ReturnType<typeof getRouter>;
	}
}
