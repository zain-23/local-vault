import type { AuditEvent } from "#/features/audit/api";
import { groupEventsByDay } from "#/features/audit/utils";
import { AuditDayGroupBlock } from "./AuditDayGroup.tsx";

type AuditTimelineProps = {
  events: AuditEvent[];
};

export function AuditTimeline({ events }: AuditTimelineProps) {
  const groups = groupEventsByDay(events);

  if (groups.length === 0) {
    return (
      <div className="rounded-lg border border-dashed border-border px-6 py-16 text-center">
        <p className="text-sm font-medium text-foreground">No activity yet</p>
        <p className="mt-1 text-sm text-muted-foreground">
          Workspace actions will show up here as your team uses LocalVault.
        </p>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-7">
      {groups.map((group) => (
        <AuditDayGroupBlock key={group.key} group={group} />
      ))}
    </div>
  );
}
