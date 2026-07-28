import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";

import {
  Button,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  Field,
  FieldError,
  FieldLabel,
  Input,
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui";
import { useInviteMember } from "#/features/members/hooks";
import {
  type InviteMemberValues,
  inviteMemberSchema,
} from "#/features/members/schemas";
import { ASSIGNABLE_ROLES } from "#/features/members/utils";
import { useModalStore } from "#/stores/useModalStore";

export function InviteMemberModal() {
  const closeModal = useModalStore((s) => s.closeModal);
  const invite = useInviteMember();
  const {
    register,
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<InviteMemberValues>({
    resolver: zodResolver(inviteMemberSchema),
    defaultValues: { email: "", role: "member" },
  });

  const onSubmit = handleSubmit((values) => {
    invite.mutate(values, {
      onSuccess: () => closeModal(),
    });
  });

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
        onSubmit={onSubmit}
        noValidate
        className="grid gap-4"
      >
        <Field data-invalid={!!errors.email}>
          <FieldLabel htmlFor="invite-email">Email address</FieldLabel>
          <Input
            id="invite-email"
            type="email"
            autoFocus
            autoComplete="email"
            placeholder="name@company.com"
            aria-invalid={!!errors.email}
            {...register("email")}
          />
          <FieldError errors={[errors.email]} />
        </Field>

        <Field data-invalid={!!errors.role}>
          <FieldLabel htmlFor="invite-role">Role</FieldLabel>
          <Controller
            control={control}
            name="role"
            render={({ field }) => (
              <Select value={field.value} onValueChange={field.onChange}>
                <SelectTrigger
                  id="invite-role"
                  className="w-full"
                  aria-invalid={!!errors.role}
                >
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
            )}
          />
          <FieldError errors={[errors.role]} />
        </Field>
      </form>

      <DialogFooter>
        <Button
          variant="outline"
          onClick={closeModal}
          disabled={invite.isPending}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          form="invite-member-form"
          isLoading={invite.isPending}
        >
          Send invite
        </Button>
      </DialogFooter>
    </DialogContent>
  );
}
