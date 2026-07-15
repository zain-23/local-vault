import type { Column, ColumnDef } from "@tanstack/react-table";
import { ArrowUpDown } from "lucide-react";
import { Avatar, AvatarFallback, AvatarImage, Button } from "#/components/ui";
import type { Member } from "#/features/members/types";
import { MemberActions } from "./MemberActions.tsx";
import { MemberRoleBadge } from "./MemberRoleBadge.tsx";

// Two-letter initials for the avatar fallback (e.g. "Ada Lovelace" -> "AL").
function initials(name: string) {
	return name
		.split(" ")
		.map((part) => part[0])
		.slice(0, 2)
		.join("")
		.toUpperCase();
}

const dateFmt = new Intl.DateTimeFormat("en-US", {
	month: "short",
	day: "numeric",
	year: "numeric",
});

// Reusable sortable header button.
function SortHeader({
	label,
	column,
}: {
	label: string;
	column: Column<Member, unknown>;
}) {
	return (
		<Button
			variant="ghost"
			size="sm"
			className="-ml-2 h-8 data-[state=open]:bg-accent"
			onClick={() => column.toggleSorting(column.getIsSorted() === "asc")}
		>
			{label}
			<ArrowUpDown className="size-3.5 opacity-60" />
		</Button>
	);
}

export const memberColumns: ColumnDef<Member>[] = [
	{
		// Identity column: avatar + name + email. Search filters against name & email.
		id: "member",
		accessorFn: (row) => row.name,
		header: "Member",
		filterFn: (row, _id, value) => {
			const q = String(value).toLowerCase();
			const m = row.original;
			return (
				m.name.toLowerCase().includes(q) || m.email.toLowerCase().includes(q)
			);
		},
		cell: ({ row }) => {
			const m = row.original;
			return (
				<div className="flex items-center gap-3">
					<Avatar>
						<AvatarImage src={m.avatarUrl} alt={m.name} />
						<AvatarFallback>{initials(m.name)}</AvatarFallback>
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
		header: ({ column }) => <SortHeader label="Role" column={column} />,
		filterFn: "equals",
		cell: ({ row }) => <MemberRoleBadge role={row.original.role} />,
	},
	{
		accessorKey: "joinedAt",
		header: ({ column }) => <SortHeader label="Joined" column={column} />,
		cell: ({ row }) => (
			<span className="text-sm text-muted-foreground">
				{dateFmt.format(new Date(row.original.joinedAt))}
			</span>
		),
	},
	{
		id: "actions",
		enableSorting: false,
		cell: ({ row }) => (
			<div className="flex justify-end">
				<MemberActions member={row.original} />
			</div>
		),
	},
];
