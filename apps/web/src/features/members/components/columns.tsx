import type { ColumnDef } from "@tanstack/react-table";

import { Avatar, AvatarImage } from "#/components/ui";
import type { Member } from "#/features/members/types";
import { dateFmt } from "../../../lib/utils.ts";
import { MemberActions } from "./MemberActions.tsx";
import { MemberRoleBadge } from "./MemberRoleBadge.tsx";

// Base roster columns (no actions). Actions are appended only for owner/admin
// so regular members never see an empty Action header.
export const memberColumns: ColumnDef<Member>[] = [
	{
		id: "member",
		accessorFn: (row) => row.name,
		header: "Member",
		cell: ({ row }) => {
			const m = row.original;
			return (
				<div className="flex items-center gap-3">
					<Avatar>
						<AvatarImage
							src={
								m.avatar_url || `https://github.com/identicons/${row.index}.png`
							}
							alt={m.name}
						/>
					</Avatar>
					<div className="min-w-0">
						<div className="truncate text-sm font-medium">{m.name}</div>
						<div className="truncate text-xs text-muted-foreground">
							{m.email}
						</div>
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
		accessorKey: "joined_at",
		header: "Joined",
		cell: ({ row }) => (
			<span className="text-sm text-muted-foreground">
				{dateFmt.format(new Date(row.original.joined_at))}
			</span>
		),
	},
];

const memberActionsColumn: ColumnDef<Member> = {
	id: "actions",
	header: "Action",
	enableSorting: false,
	cell: ({ row }) => <MemberActions member={row.original} />,
};

export function getMemberColumns(canManage: boolean): ColumnDef<Member>[] {
	return canManage ? [...memberColumns, memberActionsColumn] : memberColumns;
}
