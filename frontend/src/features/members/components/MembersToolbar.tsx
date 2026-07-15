import type { Table } from "@tanstack/react-table";
import { Search } from "lucide-react";
import { Input } from "#/components/ui";
import type { Member } from "#/features/members/types";
import { RoleFilter } from "./RoleFilter.tsx";

// Table toolbar: search on the left, role filter on the right. Search is bound to
// the "member" column's filter (matches name + email). Lives inside the table card
// so it's separate from the page's primary "Invite" action.
export function MembersToolbar({ table }: { table: Table<Member> }) {
  const memberColumn = table.getColumn("member");
  const search = (memberColumn?.getFilterValue() as string) ?? "";

  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="relative w-full sm:max-w-xs">
        <Search className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          value={search}
          onChange={(e) => memberColumn?.setFilterValue(e.target.value)}
          placeholder="Search by name or email…"
          className="pl-8"
        />
      </div>
      <RoleFilter table={table} />
    </div>
  );
}
