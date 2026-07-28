import { Search } from "lucide-react";
import { useQueryStates } from "nuqs";

import { Input } from "#/components/ui";
import {
  membersSearchDebounce,
  membersSearchOptions,
  membersSearchParams,
} from "../search-params.ts";
import { RoleFilter } from "./RoleFilter.tsx";

// Table toolbar: search + role filter, both driven by nuqs URL state so the
// list is shareable/refresh-safe and filters hit the server (not client-side).
export function MembersToolbar() {
  const [{ search }, setParams] = useQueryStates(membersSearchParams, {
    ...membersSearchOptions,
    // Only search typing needs a debounce — role/page are discrete clicks.
    limitUrlUpdates: membersSearchDebounce,
  });

  return (
    <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
      <div className="relative w-full sm:max-w-xs">
        <Search className="pointer-events-none absolute top-1/2 left-2.5 size-4 -translate-y-1/2 text-muted-foreground" />
        <Input
          value={search}
          onChange={(e) => setParams({ search: e.target.value, page: 1 })}
          placeholder="Search by name or email…"
          className="pl-8"
        />
      </div>
      <RoleFilter />
    </div>
  );
}
