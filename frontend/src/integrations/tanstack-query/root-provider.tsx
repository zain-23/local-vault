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
