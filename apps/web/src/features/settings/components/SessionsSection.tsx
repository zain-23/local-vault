import { useQuery } from "@tanstack/react-query";

import {
  Badge,
  Button,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "#/components/ui";
import { sessionsQuery } from "#/features/settings/api";
import {
  useRevokeOtherSessions,
  useRevokeSession,
} from "#/features/settings/hooks";
import {
  formatMemberSince,
  formatSessionDevice,
} from "#/features/settings/utils";

export function SessionsSection() {
  const { data: sessions = [], isLoading, isError } = useQuery(sessionsQuery);
  const revokeSession = useRevokeSession();
  const revokeOtherSessions = useRevokeOtherSessions();

  return (
    <div>
      <div className="mb-2">
        <h2 className="text-[15px] font-semibold tracking-tight">Sessions</h2>
        <p className="mt-1 text-sm text-muted-foreground">
          Devices currently signed in to your account
        </p>
      </div>

      <div className="mt-3 overflow-hidden rounded-lg border border-border">
        <Table className="min-w-140 border-collapse">
          <TableHeader>
            <TableRow className="bg-muted/40 text-left text-[11.5px] tracking-wide text-muted-foreground uppercase hover:bg-muted/40">
              <TableHead className="h-auto px-3.5 py-2.5 text-muted-foreground">
                Device
              </TableHead>
              <TableHead className="h-auto px-3.5 py-2.5 text-muted-foreground">
                IP
              </TableHead>
              <TableHead className="h-auto px-3.5 py-2.5 text-muted-foreground">
                Created
              </TableHead>
              <TableHead className="h-auto px-3.5 py-2.5 text-muted-foreground">
                Expires
              </TableHead>
              <TableHead className="h-auto w-24 px-3.5 py-2.5" />
            </TableRow>
          </TableHeader>
          <TableBody>
            {isLoading ? (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="px-3.5 py-8 text-center text-sm text-muted-foreground"
                >
                  Loading sessions...
                </TableCell>
              </TableRow>
            ) : null}
            {isError ? (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="px-3.5 py-8 text-center text-sm text-muted-foreground"
                >
                  Unable to load sessions. Please try again.
                </TableCell>
              </TableRow>
            ) : null}
            {!isLoading && !isError && sessions.length === 0 ? (
              <TableRow>
                <TableCell
                  colSpan={5}
                  className="px-3.5 py-8 text-center text-sm text-muted-foreground"
                >
                  No active sessions found.
                </TableCell>
              </TableRow>
            ) : null}
            {sessions.map((session) => (
              <TableRow key={session.id}>
                <TableCell className="px-3.5 py-3">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="font-medium">
                      {formatSessionDevice(session.user_agent)}
                    </span>
                    {session.current ? (
                      <Badge
                        variant="outline"
                        className="border-primary/40 text-primary"
                      >
                        Current
                      </Badge>
                    ) : null}
                  </div>
                </TableCell>
                <TableCell className="px-3.5 py-3 font-mono text-xs text-muted-foreground">
                  {session.ip}
                </TableCell>
                <TableCell className="px-3.5 py-3 text-[13px] text-muted-foreground">
                  {formatMemberSince(session.created_at)}
                </TableCell>
                <TableCell className="px-3.5 py-3 text-[13px] text-muted-foreground">
                  {formatMemberSince(session.expires_at)}
                </TableCell>
                <TableCell className="px-3.5 py-3 text-right">
                  {!session.current ? (
                    <Button
                      variant="destructive"
                      disabled={revokeSession.isPending}
                      onClick={() => revokeSession.mutate(session.id)}
                    >
                      Revoke
                    </Button>
                  ) : null}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {sessions.length <= 1 && !isLoading && !isError ? (
        <p className="mt-3 text-sm text-muted-foreground">
          No other sessions. Only this browser is signed in.
        </p>
      ) : (
        <div className="mt-3 flex justify-end">
          <Button
            variant="outline"
            disabled={revokeOtherSessions.isPending}
            onClick={() => revokeOtherSessions.mutate()}
          >
            {revokeOtherSessions.isPending
              ? "Signing out..."
              : "Sign out other sessions"}
          </Button>
        </div>
      )}
    </div>
  );
}
