import { useMemo } from "react";
import { useQueryStates } from "nuqs";

import { Pagination } from "#/components/shared";
import { useAuditEvents } from "#/features/audit/hooks";
import {
  auditSearchOptions,
  auditSearchParams,
  rangePresetToFrom,
} from "#/features/audit/utils";
import { useWorkspaceStore } from "#/stores";

import { AuditSkeleton } from "./AuditSkeleton.tsx";
import { AuditTimeline } from "./AuditTimeline.tsx";
import { AuditToolbar } from "./AuditToolbar.tsx";

export function AuditPage() {
  const workspace = useWorkspaceStore((s) => s.active);
  const [{ page, action_prefix, range }, setParams] = useQueryStates(
    auditSearchParams,
    auditSearchOptions,
  );

  // Memoize so `from` (derived from Date) does not change every render and
  // thrash the React Query key into a continuous refetch loop.
  const params = useMemo(
    () => ({
      page,
      limit: 40,
      action_prefix: action_prefix === "all" ? undefined : action_prefix,
      from: rangePresetToFrom(range),
    }),
    [page, action_prefix, range],
  );

  const { data, isLoading, isError, error, isFetching } =
    useAuditEvents(params);

  const events = data?.items ?? [];
  const meta = data?.meta;
  const total = meta?.total ?? 0;

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Audit log</h1>
          <p className="mt-1 text-sm text-muted-foreground">
            Every action across{" "}
            {workspace?.name ? (
              <span className="text-foreground">{workspace.name}</span>
            ) : (
              "the workspace"
            )}
            , signed and timestamped
            {total > 0 ? (
              <span className="text-muted-foreground">
                {" "}
                · {total.toLocaleString()} event{total === 1 ? "" : "s"}
              </span>
            ) : null}
          </p>
        </div>
        <AuditToolbar />
      </div>

      {isLoading ? (
        <AuditSkeleton />
      ) : isError ? (
        <div className="rounded-lg border border-destructive/30 bg-destructive/5 px-4 py-3 text-sm text-destructive">
          {error?.message || "Could not load audit events"}
        </div>
      ) : (
        <div
          className={
            isFetching ? "opacity-70 transition-opacity" : "transition-opacity"
          }
        >
          <AuditTimeline events={events} />
        </div>
      )}

      {meta ? (
        <Pagination
          page={meta.page}
          totalPages={meta.total_pages}
          onPageChange={(next) => setParams({ page: next })}
        />
      ) : null}
    </div>
  );
}
