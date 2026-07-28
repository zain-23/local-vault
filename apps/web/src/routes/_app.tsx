import { AppSidebar } from "#/components/layout/AppSidebar.tsx";
import { GlobalModalProvider } from "#/components/shared/GlobalModalProvider.tsx";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "#/components/ui";
import { meQuery } from "#/features/auth/api";
import { parseMemberRole } from "#/features/members/utils/canManageInvites.ts";
import { workspacesQuery } from "#/features/onboarding/api";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";
import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { useLayoutEffect } from "react";

export const Route = createFileRoute("/_app")({
  beforeLoad: async ({ context }) => {
    const user = await context.queryClient.ensureQueryData(meQuery);
    if (!user) {
      // stash where they were headed so login can send them back
      throw redirect({ to: "/auth/login" });
    }
    if (!user.onboarded) throw redirect({ to: "/onboarding" });

    // One fetch per page load (workspacesQuery). Do NOT
    // write Zustand here — server store mutations never reach the client.
    await context.queryClient.ensureQueryData(workspacesQuery);
  },
  component: AppLayout,
});

function AppLayout() {
  // Copy the cached workspaces result into Zustand once on the client. Store
  // stays for the rest of the session; a refresh clears both and re-fetches.
  const { data } = useQuery(workspacesQuery);
  const activeId = useWorkspaceStore((s) => s.active?.id);
  const setActive = useWorkspaceStore((s) => s.setActive);

  useLayoutEffect(() => {
    if (activeId) return;
    const first = data?.at(0);
    if (!first) return;
    setActive({
      id: first.workspace.id,
      name: first.workspace.name,
      plan: "Free",
      role: parseMemberRole(first.role),
    });
  }, [activeId, data, setActive]);

  return (
    <SidebarProvider>
      <AppSidebar />
      <SidebarInset>
        {/* Slim bar for the sidebar toggle only — full top navbar is intentionally
            skipped for now. */}
        <header className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
          <SidebarTrigger />
        </header>
        <Outlet />
      </SidebarInset>
      {/* Mounted once so any page can openModal({ type }) without rendering a Dialog. */}
      <GlobalModalProvider />
    </SidebarProvider>
  );
}
