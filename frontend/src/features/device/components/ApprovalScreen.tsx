import { CheckCircle2, Terminal } from "lucide-react";
import { useState } from "react";

import { Button, ErrorMessage, SuccessMessage } from "#/components/ui";
import { cn } from "../../../lib/utils.ts";
import { DeviceCard } from "./DeviceCard.tsx";

// UI only — no backend yet. Static request stands in for the API; Approve / Deny
// just switch the local view. Wire these to the device endpoints when the
// backend integration lands. A device links to the whole account now — the
// workspace it acts in is chosen later in the CLI, not here.
const MOCK_REQUEST = {
  deviceName: "ahmed-mbp",
  ip: "203.0.113.42",
};

// One key/value line inside the device summary panel.
function DetailRow({
  label,
  children,
}: {
  label: string;
  children: React.ReactNode;
}) {
  return (
    <div className="flex items-center justify-between gap-4 text-[13px]">
      <span className="text-muted-foreground">{label}</span>
      <span className="text-right">{children}</span>
    </div>
  );
}

function ApprovalScreen() {
  const [outcome, setOutcome] = useState<"approved" | "denied" | null>(null);

  return (
    <DeviceCard
      header={
        <div className="mb-5 flex items-start gap-3.5">
          <div className="flex size-11 shrink-0 items-center justify-center rounded-xl border border-border bg-muted text-foreground">
            <Terminal className="size-5" />
          </div>
          <div className="min-w-0">
            <h1 className="text-[17px] font-semibold tracking-[-0.01em]">
              Authorize CLI device
            </h1>
            <p className="mt-0.5 text-[13px] text-muted-foreground">
              A new terminal wants to link to your account.
            </p>
          </div>
        </div>
      }
    >
      <div className="flex flex-col gap-2.5 rounded-xl border border-border bg-muted/40 px-4 py-3.5">
        <DetailRow label="Device name">
          <span className="font-mono">{MOCK_REQUEST.deviceName}</span>
        </DetailRow>
        <DetailRow label="IP address">
          <span className="font-mono">{MOCK_REQUEST.ip}</span>
        </DetailRow>
      </div>

      <div className={cn("mt-6 flex gap-3", outcome !== null && "hidden")}>
        <Button
          variant="outline"
          className="flex-1"
          onClick={() => setOutcome("denied")}
        >
          Deny
        </Button>
        <Button
          className="flex-1"
          icon={CheckCircle2}
          onClick={() => setOutcome("approved")}
        >
          Approve
        </Button>
      </div>

      {outcome === "approved" && <SuccessMessage title="Device linked" />}
      {outcome === "denied" && <ErrorMessage title="Request denied" />}
    </DeviceCard>
  );
}

export { ApprovalScreen };
