import {
  MutationCache,
  QueryClient,
  QueryClientProvider,
} from "@tanstack/react-query";
import type { ReactNode } from "react";

export function getContext() {
  const queryClient: QueryClient = new QueryClient({
    defaultOptions: {
      queries: {
        staleTime: 60_000,
        gcTime: 5 * 60_000,
        retry: 1,
        refetchOnWindowFocus: false,
      },
      mutations: {
        retry: 0,
      },
    },
    mutationCache: new MutationCache({
      onSuccess: (_data, _var, _ctx, mutation) => {
        const keys = mutation.meta?.invalidates;
        if (!keys) return;
        // Only marks stale. Queries hydrated from SSR carry no queryFn, and the
        // ones read by route guards have no observers to borrow one from, so a
        // refetch here would reject with "Missing queryFn" and be swallowed.
        // Anything needing fresh data before navigating must fetchQuery itself.
        return Promise.all(
          keys.map((queryKey) => queryClient.invalidateQueries({ queryKey })),
        );
      },
    }),
  });

  return {
    queryClient,
  };
}

// Provides the router-context QueryClient to the React tree.
export function Provider({
  children,
  queryClient,
}: {
  children: ReactNode;
  queryClient: QueryClient;
}) {
  return (
    <QueryClientProvider client={queryClient}>{children}</QueryClientProvider>
  );
}
