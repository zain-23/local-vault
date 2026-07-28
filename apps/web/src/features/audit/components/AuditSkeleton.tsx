import { Skeleton } from "#/components/ui";

export function AuditSkeleton() {
  return (
    <div className="flex flex-col gap-7">
      {[0, 1].map((g) => (
        <div key={g} className="flex flex-col gap-2">
          <Skeleton className="h-4 w-28" />
          <div className="flex flex-col gap-3 border-t border-border pt-3">
            {[0, 1, 2, 3].map((r) => (
              <div key={r} className="flex items-center gap-3">
                <Skeleton className="h-3 w-12" />
                <Skeleton className="size-7 rounded-md" />
                <Skeleton className="h-4 flex-1" />
              </div>
            ))}
          </div>
        </div>
      ))}
    </div>
  );
}
