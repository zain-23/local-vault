import { useState } from "react";
import { toast } from "sonner";
import {
  Button,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  Field,
  FieldLabel,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui";
import { ASSIGNABLE_ROLES } from "#/features/members/types";
import { useModalStore } from "#/stores/useModalStore";

// Invite modal — fields mirror the server's InviteRequest: email + role only.
// (No "message" field; the API doesn't accept one.)
export function InviteMemberModal() {
  const closeModal = useModalStore((s) => s.closeModal);
  const [email, setEmail] = useState("");
  const [role, setRole] = useState<(typeof ASSIGNABLE_ROLES)[number]>("member");

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    // TODO: wire to POST /workspaces/:wid/members/invite. Mock for now.
    toast.success(`Invite sent to ${email}`);
    closeModal();
  }

  return (
    <DialogContent>
      <DialogHeader>
        <DialogTitle>Invite member</DialogTitle>
        <DialogDescription>
          Send an invite email. They'll join with the role you choose.
        </DialogDescription>
      </DialogHeader>

      <form
        id="invite-member-form"
        onSubmit={handleSubmit}
        className="grid gap-4"
      >
        <Field>
          <FieldLabel htmlFor="invite-email">Email address</FieldLabel>
          <Input
            id="invite-email"
            type="email"
            required
            autoFocus
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            placeholder="name@company.com"
          />
        </Field>

        <Field>
          <FieldLabel htmlFor="invite-role">Role</FieldLabel>
          <Select
            value={role}
            onValueChange={(v) =>
              setRole(v as (typeof ASSIGNABLE_ROLES)[number])
            }
          >
            <SelectTrigger id="invite-role" className="w-full">
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
        <Button variant="outline" onClick={closeModal}>
          Cancel
        </Button>
        <Button type="submit" form="invite-member-form" disabled={!email}>
          Send invite
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}
