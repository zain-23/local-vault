import { useQuery } from "@tanstack/react-query";
import { Link } from "@tanstack/react-router";
import { AlertTriangle, Check, Users } from "lucide-react";
import type * as React from "react";

import { FocusedCard } from "#/components/shared";
import { Button, StatusMessage } from "#/components/ui";
import { meQuery } from "#/features/auth/api";
import { useJoinWorkspace } from "#/features/members/hooks";
import { ApiError } from "#/services/api";

export type JoinWorkspaceSearch = {
  token?: string;
  workspaceId?: string;
};

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

function Header({ expired }: { expired: boolean }) {
  return (
    <div className="mb-5 flex items-start gap-3.5">
      <div className="flex size-11 shrink-0 items-center justify-center rounded-xl border border-border bg-muted text-foreground">
        {expired ? (
          <AlertTriangle className="size-5" />
        ) : (
          <Users className="size-5" />
        )}
      </div>
      <div className="min-w-0">
        <h1 className="text-[17px] font-semibold tracking-[-0.01em]">
          {expired ? "Invite unavailable" : "Join workspace"}
        </h1>
        <p className="mt-0.5 text-[13px] text-muted-foreground">
          {expired
            ? "This invite is invalid or has expired."
            : "You've been invited to this workspace."}
        </p>
      </div>
    </div>
  );
}

function isExpiredInviteError(error: Error | null): boolean {
  if (!error) return false;
  const msg = error.message.toLowerCase();
  return (
    msg.includes("invalid or expired") ||
    msg.includes("does not match this workspace")
  );
}

// Accept-invite screen. Mirrors the CLI authorize card: detail panel + dual CTA.
// Invite preview isn't an API yet, so we only show facts from the URL + session.
function JoinWorkspaceScreen({ token, workspaceId }: JoinWorkspaceSearch) {
  const { data: user } = useQuery(meQuery);
  const join = useJoinWorkspace();

  const linkOk = Boolean(token && workspaceId);
  const expired = !linkOk || (join.isError && isExpiredInviteError(join.error));

  const wrongEmail =
    join.error instanceof ApiError && join.error.status === 403;
  const alreadyMember =
    join.error instanceof ApiError && join.error.status === 409;

  const onAccept = () => {
    if (!token || !workspaceId) return;
    join.mutate({ workspaceId, token });
  };

  if (join.isSuccess) {
    return (
      <FocusedCard header={<Header expired={false} />}>
        <StatusMessage.Success title="You're in — heading to the workspace." />
      </FocusedCard>
    );
  }

  if (expired) {
    return (
      <FocusedCard header={<Header expired />}>
        {linkOk && (
          <div className="flex flex-col gap-2.5 rounded-xl border border-border bg-muted/40 px-4 py-3.5 opacity-70">
            <DetailRow label="Workspace">
              <span className="font-mono text-[12px] text-muted-foreground">
                {workspaceId}
              </span>
            </DetailRow>
          </div>
        )}
        <StatusMessage.WarnBanner className="mt-3.5">
          Ask a workspace owner or admin to send a new invite.
        </StatusMessage.WarnBanner>
        <Button asChild variant="outline" className="mt-4.5 w-full">
          <Link to="/">{user?.onboarded ? "Back to dashboard" : "Back"}</Link>
        </Button>
      </FocusedCard>
    );
  }

  return (
    <FocusedCard header={<Header expired={false} />}>
      <div className="flex flex-col gap-2.5 rounded-xl border border-border bg-muted/40 px-4 py-3.5">
        <DetailRow label="Workspace">
          <span className="font-mono text-[12px]">{workspaceId}</span>
        </DetailRow>
        {user?.email && (
          <DetailRow label="Invited as">
            <span className="font-mono text-[12px]">{user.email}</span>
          </DetailRow>
        )}
      </div>

      <StatusMessage.InfoBanner className="mt-3.5">
        You'll join with the role on the invite. You can leave later from
        workspace settings.
      </StatusMessage.InfoBanner>

      {(wrongEmail || alreadyMember || (join.isError && !expired)) && (
        <StatusMessage.Error
          className="mt-3.5"
          title={
            wrongEmail
              ? "This invite was sent to a different email."
              : alreadyMember
                ? "You're already a member of this workspace."
                : (join.error?.message ?? "Could not join workspace.")
          }
        />
      )}

      <div className="mt-4.5 flex gap-3">
        <Button asChild variant="outline" className="flex-1">
          <Link to="/">Decline</Link>
        </Button>
        <Button
          className="flex-1"
          icon={Check}
          isLoading={join.isPending}
          onClick={onAccept}
        >
          Accept invite
        </Button>
      </div>
    </FocusedCard>
  );
}

export { JoinWorkspaceScreen };
