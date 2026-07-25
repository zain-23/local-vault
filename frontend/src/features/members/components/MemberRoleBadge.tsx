import { Badge } from "#/components/ui";
import type { MemberRole } from "#/features/members/types";

// Owner stands out (primary), admin is neutral-strong (secondary), member is quiet (outline).
const ROLE_VARIANT: Record<
  MemberRole,
  React.ComponentProps<typeof Badge>["variant"]
> = {
  owner: "default",
  admin: "secondary",
  member: "outline",
};

export function MemberRoleBadge({ role }: { role: MemberRole }) {
  return (
    <Badge variant={ROLE_VARIANT[role]} className="capitalize">
      {role}
    </Badge>
  );
}
