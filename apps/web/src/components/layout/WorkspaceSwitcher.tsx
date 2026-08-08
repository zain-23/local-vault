import { useQuery } from "@tanstack/react-query";
import { Check, ChevronsUpDown, Plus } from "lucide-react";
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
} from "#/components/ui";
import { parseMemberRole } from "#/features/members/utils/canManageInvites.ts";
import { workspacesQuery } from "#/features/onboarding/api";
import { WorkspaceIcon } from "#/features/onboarding/lib/workspaceIcons.tsx";
import { useModalStore } from "#/stores/useModalStore";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";

export function WorkspaceSwitcher() {
	const { data: workspaces = [] } = useQuery(workspacesQuery);
	const workspace = useWorkspaceStore((s) => s.active);
	const setActive = useWorkspaceStore((s) => s.setActive);
	const openModal = useModalStore((s) => s.openModal);

	return (
		<SidebarMenu>
			<SidebarMenuItem>
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<SidebarMenuButton
							size="lg"
							className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
						>
							<WorkspaceIcon id={workspace?.icon} size="md" />
							<div className="grid flex-1 text-left leading-tight">
								<span className="truncate font-medium">
									{workspace?.name ?? "Workspace"}
								</span>
							</div>
							<ChevronsUpDown className="ml-auto size-4" />
						</SidebarMenuButton>
					</DropdownMenuTrigger>
					<DropdownMenuContent
						className="w-(--radix-dropdown-menu-trigger-width)"
						align="start"
					>
						{workspaces.map(({ workspace: ws, role }) => {
							const selected = ws.id === workspace?.id;
							return (
								<DropdownMenuItem
									key={ws.id}
									onSelect={() =>
										setActive({
											id: ws.id,
											name: ws.name,
											icon: ws.icon,
											role: parseMemberRole(role),
										})
									}
								>
									<WorkspaceIcon id={ws.icon} size="sm" />
									<span className="truncate">{ws.name}</span>
									{selected && <Check className="ml-auto size-4" />}
								</DropdownMenuItem>
							);
						})}
						<DropdownMenuSeparator />
						<DropdownMenuItem
							onSelect={() => openModal({ type: "create-workspace" })}
						>
							<Plus />
							<span>Create workspace</span>
						</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</SidebarMenuItem>
		</SidebarMenu>
	);
}
