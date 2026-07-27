import {
  type ColumnDef,
  type ColumnFiltersState,
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  type SortingState,
  type Table as TableInstance,
  useReactTable,
} from "@tanstack/react-table";
import { useState } from "react";
import { Spinner } from "./Spinner";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./Table";

interface DataTableProps<TData, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  // Optional toolbar (search, filters, …) rendered above the table. Receives the
  // table instance so controls can drive sorting/filtering generically.
  toolbar?: (table: TableInstance<TData>) => React.ReactNode;
  // Shown when there are no rows (after filtering), and not loading/errored.
  emptyMessage?: string;
  isLoading?: boolean;
  // When set (and not loading), replaces the empty/data body with this message.
  errorMessage?: string;
  // Optional row click (e.g. navigate to detail). Keyboard: Enter/Space on row.
  onRowClick?: (row: TData) => void;
}

// Generic, reusable table: owns TanStack sorting + column-filter state and renders
// with the shadcn Table primitives. Feature code only supplies columns + data.
export function DataTable<TData, TValue>({
  columns,
  data,
  toolbar,
  emptyMessage = "No results.",
  isLoading = false,
  errorMessage,
  onRowClick,
}: DataTableProps<TData, TValue>) {
  const [sorting, setSorting] = useState<SortingState>([]);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const table = useReactTable({
    data,
    columns,
    state: { sorting, columnFilters },
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });

  const rows = table.getRowModel().rows;
  const showStatus = isLoading || !!errorMessage || rows.length === 0;

  return (
    <div className="flex flex-col gap-4">
      {toolbar?.(table)}
      <div className="overflow-hidden rounded-md border">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id}>
                {headerGroup.headers.map((header) => (
                  <TableHead
                    key={header.id}
                    style={{ width: header.getSize() }}
                  >
                    {header.isPlaceholder
                      ? null
                      : flexRender(
                          header.column.columnDef.header,
                          header.getContext(),
                        )}
                  </TableHead>
                ))}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody>
            {showStatus ? (
              <TableRow>
                <TableCell
                  colSpan={columns.length}
                  className="h-24 text-center text-muted-foreground"
                >
                  {isLoading ? (
                    <span className="inline-flex items-center justify-center gap-2">
                      <Spinner size="default" />
                      <span className="sr-only">Loading</span>
                    </span>
                  ) : errorMessage ? (
                    errorMessage
                  ) : (
                    emptyMessage
                  )}
                </TableCell>
              </TableRow>
            ) : (
              rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                  className={onRowClick ? "cursor-pointer" : undefined}
                  tabIndex={onRowClick ? 0 : undefined}
                  onClick={
                    onRowClick ? () => onRowClick(row.original) : undefined
                  }
                  onKeyDown={
                    onRowClick
                      ? (e) => {
                          if (e.key === "Enter" || e.key === " ") {
                            e.preventDefault();
                            onRowClick(row.original);
                          }
                        }
                      : undefined
                  }
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>
                      {flexRender(
                        cell.column.columnDef.cell,
                        cell.getContext(),
                      )}
                    </TableCell>
                  ))}
                </TableRow>
              ))
            )}
          </TableBody>
        </Table>
      </div>
    </div>
  );
}
