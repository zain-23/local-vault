import { Link, useRouterState } from "@tanstack/react-router";
import {
	LayoutDashboard,
	type LucideIcon,
	ScrollText,
	Settings,
	Users,
	Vault,
} from "lucide-react";
import {
	Sidebar,
	SidebarContent,
	SidebarFooter,
	SidebarGroup,
	SidebarGroupLabel,
	SidebarHeader,
	SidebarMenu,
	SidebarMenuBadge,
	SidebarMenuButton,
	SidebarMenuItem,
	SidebarRail,
} from "#/components/ui";
import { NavUser } from "./NavUser.tsx";
import { WorkspaceSwitcher } from "./WorkspaceSwitcher.tsx";

interface NavItem {
	title: string;
	icon: LucideIcon;
	to?: string;
}

const NAV_ITEMS: NavItem[] = [
	{ title: "Dashboard", icon: LayoutDashboard, to: "/dashboard" },
	{ title: "Vaults", icon: Vault, to: "/vaults" },
	{ title: "Members", icon: Users, to: "/members" },
	{ title: "Audit log", icon: ScrollText, to: "/audit" },
	{ title: "Settings", icon: Settings, to: "/settings" },
];

export function AppSidebar() {
	const pathname = useRouterState({
		select: (s) => s.location.pathname,
	});

	return (
		<Sidebar collapsible="icon">
			<SidebarHeader>
				<WorkspaceSwitcher />
			</SidebarHeader>

			<SidebarContent>
				<SidebarGroup>
					<SidebarGroupLabel>Workspace</SidebarGroupLabel>
					<SidebarMenu>
						{NAV_ITEMS.map((item) => (
							<SidebarMenuItem key={item.title}>
								{item.to ? (
									<SidebarMenuButton
										asChild
										tooltip={item.title}
										isActive={
											pathname === item.to ||
											(item.to !== undefined &&
												pathname.startsWith(`${item.to}/`))
										}
									>
										<Link to={item.to}>
											<item.icon />
											<span>{item.title}</span>
										</Link>
									</SidebarMenuButton>
								) : (
									<>
										<SidebarMenuButton tooltip={item.title} disabled>
											<item.icon />
											<span>{item.title}</span>
										</SidebarMenuButton>
										<SidebarMenuBadge>Soon</SidebarMenuBadge>
									</>
								)}
							</SidebarMenuItem>
						))}
					</SidebarMenu>
				</SidebarGroup>
			</SidebarContent>

			<SidebarFooter>
				<NavUser />
			</SidebarFooter>

			<SidebarRail />
		</Sidebar>
	);
}
