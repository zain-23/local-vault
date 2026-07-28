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
import { useModalStore } from "#/stores/useModalStore";

// Row action menu. Only exposes operations the member API actually supports:
// change role (PUT /:userId/role) and remove (DELETE /:userId). The owner row
// gets no actions — the server forbids changing or removing the owner.
export function MemberActions({ member }: { member: Member }) {
  const openModal = useModalStore((s) => s.openModal);

  if (member.role === "owner") {
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
        <DropdownMenuItem
          onSelect={() => openModal({ type: "change-role", props: { member } })}
        >
          <Shield />
          Change role
        </DropdownMenuItem>
        <DropdownMenuSeparator />
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
