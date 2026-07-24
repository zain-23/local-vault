import { ChevronsUpDown } from "lucide-react";
import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from "#/components/ui";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";

// Sidebar header showing the active workspace identity. Switching between multiple
// workspaces isn't wired yet (only one exists post-onboarding), so this is a
// single-workspace display for now — structured to grow into a real switcher.
export function WorkspaceSwitcher() {
  const workspace = useWorkspaceStore((s) => s.active);
  const initial = workspace?.name?.charAt(0).toUpperCase() ?? "W";

  return (
    <SidebarMenu>
      <SidebarMenuItem>
        <SidebarMenuButton
          size={"lg"}
          className="data-[state=open]:bg-sidebar-accent"
        >
          <div className="flex aspect-square size-8 items-center justify-center rounded-md bg-sidebar-primary text-sm font-semibold text-sidebar-primary-foreground">
            {initial}
          </div>
          <div className="grid flex-1 text-left leading-tight">
            <span className="truncate font-medium">{workspace?.name}</span>
            <span className="truncate text-xs text-muted-foreground">
              {workspace?.plan} plan
            </span>
          </div>
          <ChevronsUpDown className="ml-auto size-4 opacity-50" />
        </SidebarMenuButton>
      </SidebarMenuItem>
    </SidebarMenu>
  );
}
