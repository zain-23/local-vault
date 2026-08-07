import { useQuery } from "@tanstack/react-query";
import { createFileRoute, Outlet } from "@tanstack/react-router";
import { useLayoutEffect } from "react";
import { AppBreadcrumb } from "#/components/layout/AppBreadcrumb.tsx";
import { AppSidebar } from "#/components/layout/AppSidebar.tsx";
import { GlobalModalProvider } from "#/components/shared/GlobalModalProvider.tsx";
import { SidebarInset, SidebarProvider, SidebarTrigger } from "#/components/ui";
import { authGuard } from "#/features/auth/utils/authGuard";
import { parseMemberRole } from "#/features/members/utils/canManageInvites.ts";
import { workspacesQuery } from "#/features/onboarding/api";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";

export const Route = createFileRoute("/_app")({
	...authGuard(
		(user) => {
			if (!user) return { to: "/auth/login" };
			if (!user.onboarded) return { to: "/onboarding" };
		},
		// write Zustand here — server store mutations never reach the client.
		(context) => context.queryClient.ensureQueryData(workspacesQuery),
	),
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
				<header className="flex h-12 shrink-0 items-center gap-2 border-b px-4">
					<SidebarTrigger />
					<AppBreadcrumb />
				</header>
				<Outlet />
			</SidebarInset>
			{/* Mounted once so any page can openModal({ type }) without rendering a Dialog. */}
			<GlobalModalProvider />
		</SidebarProvider>
	);
}
