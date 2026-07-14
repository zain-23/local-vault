import { createFileRoute } from "@tanstack/react-router";

import { LoginForm } from "#/features/auth/components/index.ts";

export const Route = createFileRoute("/auth/login")({
  validateSearch: (search: Record<string, unknown>): { error?: string } => {
    return typeof search.error === "string" ? { error: search.error } : {};
  },
  component: LoginForm,
});
