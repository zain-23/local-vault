import type { ColumnDef } from "@tanstack/react-table";
import { Vault } from "lucide-react";

import { Badge } from "#/components/ui";
import type { VaultSummary } from "#/features/vaults/api";
import { dateFmt } from "#/lib/utils.ts";

export const vaultColumns: ColumnDef<VaultSummary>[] = [
  {
    accessorKey: "name",
    header: "Name",
    cell: ({ row }) => (
      <div className="flex items-center gap-2.5">
        <Vault className="size-4 shrink-0 text-muted-foreground" />
        <span className="font-mono text-[13.5px] font-medium tracking-tight">
          {row.original.name}
        </span>
      </div>
    ),
  },
  {
    accessorKey: "peer_count",
    header: "Devices",
    size: 100,
    cell: ({ row }) => (
      <span className="font-mono text-sm text-muted-foreground">
        {row.original.peer_count}
      </span>
    ),
  },
  {
    accessorKey: "has_snapshot",
    header: "Snapshot",
    size: 120,
    cell: ({ row }) =>
      row.original.has_snapshot ? (
        <Badge variant="secondary">Yes</Badge>
      ) : (
        <Badge variant="outline">No</Badge>
      ),
  },
  {
    accessorKey: "updated_at",
    header: "Updated",
    size: 160,
    cell: ({ row }) => (
      <span className="text-sm text-muted-foreground">
        {dateFmt.format(new Date(row.original.updated_at))}
      </span>
    ),
  },
];
