import { useState } from "react";

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
import { MOCK_SESSIONS, type MockSession } from "#/features/settings/mock.ts";

export function SessionsSection() {
  const [sessions, setSessions] = useState<MockSession[]>(MOCK_SESSIONS);

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
            {sessions.map((session) => (
              <TableRow key={session.id}>
                <TableCell className="px-3.5 py-3">
                  <div className="flex flex-wrap items-center gap-2">
                    <span className="font-medium">{session.device}</span>
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
                  {session.createdAt}
                </TableCell>
                <TableCell className="px-3.5 py-3 text-[13px] text-muted-foreground">
                  {session.expiresAt}
                </TableCell>
                <TableCell className="px-3.5 py-3 text-right">
                  {!session.current ? (
                    <Button variant="destructive">Revoke</Button>
                  ) : null}
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </div>

      {sessions.length <= 1 ? (
        <p className="mt-3 text-sm text-muted-foreground">
          No other sessions. Only this browser is signed in.
        </p>
      ) : null}
    </div>
  );
}
