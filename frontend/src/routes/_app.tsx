import { createFileRoute, Outlet } from "@tanstack/react-router";
import { AppSidebar } from "#/components/layout/AppSidebar.tsx";
import { GlobalModalProvider } from "#/components/shared/GlobalModalProvider.tsx";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "#/components/ui";

// Protected app shell: sidebar + content outlet + the global modal host.
// TODO(auth): add a beforeLoad guard that redirects unauthenticated users to
// /auth/login once the frontend has a session/me query to check against.
export const Route = createFileRoute("/_app")({
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
