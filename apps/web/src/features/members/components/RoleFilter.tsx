import { useQueryStates } from "nuqs";

import {
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "#/components/ui";
import type { MemberRole } from "#/features/members/types";
import {
	ROLE_FILTER_ALL,
	ROLE_FILTER_LABELS,
	type RoleFilterValue,
} from "#/features/members/utils";
import { membersSearchOptions, membersSearchParams } from "../search-params.ts";

// Role filter backed by the URL (`?role=`). Changing role resets to page 1 so
// you never land on an empty page after narrowing the filter.
export function RoleFilter() {
	const [{ role }, setParams] = useQueryStates(
		membersSearchParams,
		membersSearchOptions,
	);

	const value: RoleFilterValue = role ?? ROLE_FILTER_ALL;

	return (
		<Select
			value={value}
			onValueChange={(next) =>
				setParams({
					role: next === ROLE_FILTER_ALL ? null : (next as MemberRole),
					page: 1,
				})
			}
		>
			<SelectTrigger className="w-36" aria-label="Filter by role">
				<SelectValue />
			</SelectTrigger>
			<SelectContent>
				{[...ROLE_FILTER_LABELS].map(([filterValue, label]) => (
					<SelectItem key={filterValue} value={filterValue}>
						{label}
					</SelectItem>
				))}
			</SelectContent>
		</Select>
	);
}
