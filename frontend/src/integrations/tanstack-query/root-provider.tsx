import {
  QueryClient,
  type QueryClientConfig,
  QueryClientProvider,
} from "@tanstack/react-query";
import type { ReactNode } from "react";

// App-wide React Query defaults.
const queryConfig: QueryClientConfig = {
  defaultOptions: {
    queries: {
      staleTime: 60_000, // 1 min — data is "fresh" this long, no refetch on mount
      gcTime: 5 * 60_000, // keep unused cache 5 min before garbage-collecting
      retry: 1, // retry a failed query once
      refetchOnWindowFocus: false, // don't refetch every time the tab refocuses
    },
    mutations: {
      retry: 0, // never auto-retry mutations — they have side effects
    },
  },
};

// A fresh client per request (getContext runs in getRouter) so server-rendered
// cache never leaks between users.
export function getContext() {
  const queryClient = new QueryClient(queryConfig);

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
