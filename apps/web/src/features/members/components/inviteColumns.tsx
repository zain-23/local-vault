import type { ColumnDef } from "@tanstack/react-table";

import { Avatar, AvatarImage } from "#/components/ui";
import type { Invite } from "#/features/members/types";
import { dateFmt } from "../../../lib/utils.ts";
import { InviteActions } from "./InviteActions.tsx";
import { MemberRoleBadge } from "./MemberRoleBadge.tsx";

export const inviteColumns: ColumnDef<Invite>[] = [
	{
		id: "invite",
		accessorKey: "email",
		header: "Invite",
		cell: ({ row }) => {
			const invite = row.original;
			return (
				<div className="flex items-center gap-3">
					<Avatar>
						<AvatarImage
							src={`https://github.com/identicons/${row.index}.png`}
							alt={invite.email}
						/>
					</Avatar>
					<div className="min-w-0">
						<div className="truncate text-sm font-medium">{invite.email}</div>
					</div>
				</div>
			);
		},
	},
	{
		accessorKey: "role",
		header: "Role",
		cell: ({ row }) => <MemberRoleBadge role={row.original.role} />,
	},
	{
		accessorKey: "created_at",
		header: "Invited",
		cell: ({ row }) => (
			<span className="text-sm text-muted-foreground">
				{dateFmt.format(new Date(row.original.created_at))}
			</span>
		),
	},
	{
		accessorKey: "expires_at",
		header: "Expires",
		cell: ({ row }) => (
			<span className="text-sm text-muted-foreground">
				{dateFmt.format(new Date(row.original.expires_at))}
			</span>
		),
	},
	{
		id: "actions",
		header: "Action",
		enableSorting: false,
		cell: ({ row }) => <InviteActions invite={row.original} />,
	},
];
