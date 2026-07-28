import type { Table } from "@tanstack/react-table";
import type { Member, MemberRole } from "#/features/members/types";
import { cn } from "#/lib/utils.ts";

const FILTERS: { label: string; value: MemberRole | "all" }[] = [
	{ label: "All", value: "all" },
	{ label: "Owner", value: "owner" },
	{ label: "Admin", value: "admin" },
	{ label: "Member", value: "member" },
];

// Segmented role filter with live counts. Drives the table's "role" column filter
// so filtering stays inside TanStack (no ad-hoc array slicing).
export function RoleFilter({ table }: { table: Table<Member> }) {
	const roleColumn = table.getColumn("role");
	const active =
		(roleColumn?.getFilterValue() as MemberRole | undefined) ?? "all";

	// Counts come from the unfiltered rows so each chip shows its true total.
	const rows = table.getCoreRowModel().rows;
	const countFor = (value: MemberRole | "all") =>
		value === "all"
			? rows.length
			: rows.filter((r) => r.original.role === value).length;

	return (
		<div className="inline-flex items-center gap-1 rounded-md border p-0.5">
			{FILTERS.map(({ label, value }) => {
				const isActive = active === value;
				return (
					<button
						key={value}
						type="button"
						onClick={() =>
							roleColumn?.setFilterValue(value === "all" ? undefined : value)
						}
						className={cn(
							"inline-flex items-center gap-1.5 rounded px-2.5 py-1 text-xs font-medium transition-colors",
							isActive
								? "bg-muted text-foreground"
								: "text-muted-foreground hover:text-foreground",
						)}
					>
						{label}
						<span className="font-mono text-[11px] text-muted-foreground/80">
							{countFor(value)}
						</span>
					</button>
				);
			})}
		</div>
	);
}
