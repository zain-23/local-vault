import type { AuditDayGroup } from "#/features/audit/utils";
import { AuditEventRow } from "./AuditEventRow.tsx";

type AuditDayGroupProps = {
  group: AuditDayGroup;
};

export function AuditDayGroupBlock({ group }: AuditDayGroupProps) {
  return (
    <section>
      <h3 className="mb-2.5 border-b border-border py-1.5 text-[11px] font-semibold tracking-[0.06em] text-muted-foreground uppercase">
        {group.label}
      </h3>
      <div className="flex flex-col">
        {group.events.map((event) => (
          <AuditEventRow key={event.id} event={event} />
        ))}
      </div>
    </section>
  );
}
