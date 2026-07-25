import { useState } from "react";
import {
  Button,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  Field,
  FieldLabel,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui";
import { useChangeRole } from "#/features/members/hooks";
import type { Member } from "#/features/members/types";
import { ASSIGNABLE_ROLES } from "#/features/members/utils";
import { useModalStore } from "#/stores/useModalStore";

// Change-role modal — PUT /workspaces/:wid/members/:userId/role.
export function ChangeRoleModal() {
  const { props, closeModal } = useModalStore();
  const member = props.member as Member | undefined;
  const changeRole = useChangeRole();
  const initial = member?.role === "admin" ? "admin" : "member";
  const [role, setRole] = useState<(typeof ASSIGNABLE_ROLES)[number]>(initial);

  if (!member) return null;

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    changeRole.mutate(
      { userId: member!.user_id, role },
      {
        onSuccess: () => closeModal(),
      },
    );
  }

  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Change role</DialogTitle>
        <DialogDescription>
          Update the workspace role for {member.name}.
        </DialogDescription>
      </DialogHeader>

      <form
        id="change-role-form"
        onSubmit={handleSubmit}
        className="grid gap-4"
      >
        <Field>
          <FieldLabel htmlFor="change-role-select">Role</FieldLabel>
          <Select
            value={role}
            onValueChange={(v) =>
              setRole(v as (typeof ASSIGNABLE_ROLES)[number])
            }
          >
            <SelectTrigger id="change-role-select" className="w-full">
              <SelectValue />
            </SelectTrigger>
            <SelectContent>
              {ASSIGNABLE_ROLES.map((r) => (
                <SelectItem key={r} value={r} className="capitalize">
                  {r}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </Field>
      </form>

      <DialogFooter>
        <Button
          variant="outline"
          onClick={closeModal}
          disabled={changeRole.isPending}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          form="change-role-form"
          disabled={role === member.role}
          isLoading={changeRole.isPending}
        >
          Save changes
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}
