import { useQuery } from "@tanstack/react-query";
import { ChevronsUpDown, LogOut } from "lucide-react";

import {
	Avatar,
	AvatarFallback,
	AvatarImage,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuTrigger,
	SidebarMenu,
	SidebarMenuButton,
	SidebarMenuItem,
	useSidebar,
} from "#/components/ui";
import { meQuery } from "#/features/auth/api";
import { useLogout } from "#/features/auth/hooks";
import { initialsFromName } from "#/features/settings/utils";

export function NavUser() {
	const { isMobile } = useSidebar();
	const { data: user } = useQuery(meQuery);
	const logout = useLogout();

	if (!user) return null;

	const initials = initialsFromName(user.name || user.email || "?");

	return (
		<SidebarMenu>
			<SidebarMenuItem>
				<DropdownMenu>
					<DropdownMenuTrigger asChild>
						<SidebarMenuButton
							size="lg"
							className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
						>
							<Avatar size="sm" className="rounded-md">
								{user.avatar_url ? (
									<AvatarImage src={user.avatar_url} alt={user.name} />
								) : null}
								<AvatarFallback className="rounded-md bg-sidebar-primary text-sidebar-primary-foreground">
									{initials}
								</AvatarFallback>
							</Avatar>
							<div className="grid flex-1 text-left leading-tight">
								<span className="truncate font-medium">{user.name}</span>
								<span className="truncate text-xs text-muted-foreground">
									{user.email}
								</span>
							</div>
							<ChevronsUpDown className="ml-auto size-4" />
						</SidebarMenuButton>
					</DropdownMenuTrigger>
					<DropdownMenuContent
						className="w-(--radix-dropdown-menu-trigger-width) min-w-56"
						side={isMobile ? "bottom" : "top"}
						align="start"
						sideOffset={4}
					>
						<DropdownMenuItem
							disabled={logout.isPending}
							onSelect={() => logout.mutate()}
						>
							<LogOut />
							Log out
						</DropdownMenuItem>
					</DropdownMenuContent>
				</DropdownMenu>
			</SidebarMenuItem>
		</SidebarMenu>
	);
}
