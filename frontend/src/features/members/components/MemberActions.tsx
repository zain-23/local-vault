import { MoreHorizontal, Shield, Trash2 } from "lucide-react";
import {
  Button,
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "#/components/ui";
import type { Member } from "#/features/members/types";
import {
  canChangeMemberRole,
  canManageInvites,
} from "#/features/members/utils";
import { useModalStore } from "#/stores/useModalStore";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";

type MemberActionsProps = {
  member: Member;
};

// Row action menu. Owner/admin only, never on the owner row. Mutations live in
// ChangeRoleModal / RemoveMemberModal — this only opens them.
export function MemberActions({ member }: MemberActionsProps) {
  const myRole = useWorkspaceStore((s) => s.active?.role);
  const openModal = useModalStore((s) => s.openModal);

  const canManage = canManageInvites(myRole);
  const canChangeRole = canChangeMemberRole(myRole);

  if (!canManage || member.role === "owner") {
    return null;
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="ghost" size="icon-sm" aria-label="Member actions">
          <MoreHorizontal />
        </Button>
      </DropdownMenuTrigger>
      <DropdownMenuContent align="end" className="w-44">
        <DropdownMenuLabel>Actions</DropdownMenuLabel>
        {canChangeRole && (
          <DropdownMenuItem
            onSelect={() =>
              openModal({ type: "change-role", props: { member } })
            }
          >
            <Shield />
            Change role
          </DropdownMenuItem>
        )}
        {canChangeRole && <DropdownMenuSeparator />}
        <DropdownMenuItem
          variant="destructive"
          onSelect={() =>
            openModal({ type: "remove-member", props: { member } })
          }
        >
          <Trash2 />
          Remove member
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}
