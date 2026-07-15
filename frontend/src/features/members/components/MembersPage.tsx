import { UserPlus } from "lucide-react";
import { Button, DataTable } from "#/components/ui";
import { MOCK_MEMBERS } from "#/features/members/mock/members.ts";
import { useModalStore } from "#/stores/useModalStore";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";
import { memberColumns } from "./columns.tsx";
import { MembersToolbar } from "./MembersToolbar.tsx";

// Members page: header with the primary "Invite" action, then the reusable
// DataTable driven by its own search + role filter toolbar. Data is mock for now;
// swapping in a react-query hook here is the only change needed to go live.
export function MembersPage() {
  const workspace = useWorkspaceStore((s) => s.active);
  const openModal = useModalStore((s) => s.openModal);
  const members = MOCK_MEMBERS;

  return (
    <div className="flex flex-col gap-6 p-6">
      <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 className="text-xl font-semibold tracking-tight">Members</h1>
          <p className="text-sm text-muted-foreground">
            {members.length} people in {workspace?.name}
          </p>
        </div>
        <Button
          icon={UserPlus}
          onClick={() => openModal({ type: "invite-member" })}
        >
          Invite member
        </Button>
      </div>

      <DataTable
        columns={memberColumns}
        data={members}
        toolbar={(table) => <MembersToolbar table={table} />}
        emptyMessage="No members match your filters."
      />
    </div>
  );
}
