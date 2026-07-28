import { Check, ChevronsUpDown } from "lucide-react";
import { useQuery } from "@tanstack/react-query";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "#/components/ui";
import { workspacesQuery } from "#/features/onboarding/api";
import { parseMemberRole } from "#/features/members/utils/canManageInvites.ts";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";

export function WorkspaceSwitcher() {
  const { data: workspaces = [] } = useQuery(workspacesQuery);
  const workspace = useWorkspaceStore((s) => s.active);
  const setActive = useWorkspaceStore((s) => s.setActive);

  const initial = workspace?.name?.charAt(0).toUpperCase() ?? "W";

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
            >
              <div className="flex aspect-square size-8 items-center justify-center rounded-md bg-sidebar-primary text-sm font-semibold text-sidebar-primary-foreground">
                {initial}
              </div>
              <div className="grid flex-1 text-left leading-tight">
                <span className="truncate font-medium">
                  {workspace?.name ?? "Workspace"}
                </span>
                <span className="truncate text-xs text-muted-foreground">
                  {workspace?.plan ?? "Free"} plan
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
                      plan: "Free",
                      role: parseMemberRole(role),
                    })
                  }
                >
                  <span className="truncate">{ws.name}</span>
                  {selected && <Check className="ml-auto size-4" />}
                </DropdownMenuItem>
              );
            })}
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
