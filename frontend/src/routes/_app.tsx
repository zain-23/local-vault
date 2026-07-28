import { createFileRoute, Outlet, redirect } from "@tanstack/react-router";
import { AppSidebar } from "#/components/layout/AppSidebar.tsx";
import { GlobalModalProvider } from "#/components/shared/GlobalModalProvider.tsx";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "#/components/ui";
import { meQuery } from "#/features/auth/api";

export const Route = createFileRoute("/_app")({
  beforeLoad: async ({ context }) => {
    const user = await context.queryClient.ensureQueryData(meQuery);
    if (!user) {
      // stash where they were headed so login can send them back
      throw redirect({ to: "/auth/login" });
    }
    if (!user.onboarded) throw redirect({ to: "/onboarding" });
  },
  component: AppLayout,
});

function AppLayout() {
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
