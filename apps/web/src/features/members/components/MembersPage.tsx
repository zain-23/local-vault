import { useQueryStates } from "nuqs";
import { useEffect } from "react";

import { Pagination, TabGroup } from "#/components/shared";
import { Button, DataTable } from "#/components/ui";
import type { PageMeta } from "#/features/members/api";
import { useInvites, useMembers } from "#/features/members/hooks";
import type { Invite, Member } from "#/features/members/types";
import { canManageInvites } from "#/features/members/utils";
import { useModalStore } from "#/stores/useModalStore";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";
import {
	type MemberTab,
	membersSearchOptions,
	membersSearchParams,
} from "../search-params.ts";
import { getMemberColumns } from "./columns.tsx";
import { inviteColumns } from "./inviteColumns.tsx";
import { MembersToolbar } from "./MembersToolbar.tsx";

// Members page: one table at a time. Owner/admin get Members | Pending invites
// tabs; regular members only see the roster. Filters stay in the URL via nuqs.
export function MembersPage() {
	const workspace = useWorkspaceStore((s) => s.active);
	const openModal = useModalStore((s) => s.openModal);
	const canManage = canManageInvites(workspace?.role);

	const [{ search, role, page, tab }, setParams] = useQueryStates(
		membersSearchParams,
		membersSearchOptions,
	);

	// Members can't open the invites tab — snap URL back if somehow set.
	useEffect(() => {
		if (!canManage && tab === "invites") {
			void setParams({ tab: "members" });
		}
	}, [canManage, tab, setParams]);

	const activeTab: MemberTab = canManage ? tab : "members";

	const membersQuery = useMembers({
		page,
		search: search || undefined,
		role: role ?? undefined,
	});
	const invitesQuery = useInvites();

	const members = membersQuery.data?.items ?? [];
	const meta = membersQuery.data?.meta;
	const invites = invitesQuery.data ?? [];
	const memberTotal = meta?.total ?? members.length;
	const inviteTotal = invites.length;

	return (
		<div className="flex flex-col gap-6 p-6">
			<div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
				<div>
					<h1 className="text-xl font-semibold tracking-tight">Members</h1>
					<p className="text-sm text-muted-foreground">
						Members and Pending invites are show here
					</p>
				</div>
				{canManage && (
					<Button onClick={() => openModal({ type: "invite-member" })}>
						Invite member
					</Button>
				)}
			</div>

			{canManage ? (
				<TabGroup
					value={activeTab}
					onValueChange={(value) =>
						setParams({ tab: value as MemberTab, page: 1 })
					}
					items={[
						{
							value: "members",
							label: "Members",
							count: memberTotal,
							content: (
								<MembersTable
									members={members}
									canManage={canManage}
									isPending={membersQuery.isPending}
									isError={membersQuery.isError}
									errorMessage={membersQuery.error?.message}
									meta={meta}
									onPageChange={(next) => setParams({ page: next })}
								/>
							),
						},
						{
							value: "invites",
							label: "Pending invites",
							count: inviteTotal,
							content: (
								<InvitesTable
									invites={invites}
									isPending={invitesQuery.isPending}
									isError={invitesQuery.isError}
									errorMessage={invitesQuery.error?.message}
								/>
							),
						},
					]}
				/>
			) : (
				<MembersTable
					members={members}
					canManage={false}
					isPending={membersQuery.isPending}
					isError={membersQuery.isError}
					errorMessage={membersQuery.error?.message}
					meta={meta}
					onPageChange={(next) => setParams({ page: next })}
				/>
			)}
		</div>
	);
}

function MembersTable({
	members,
	canManage,
	isPending,
	isError,
	errorMessage,
	meta,
	onPageChange,
}: {
	members: Member[];
	canManage: boolean;
	isPending: boolean;
	isError: boolean;
	errorMessage?: string;
	meta: PageMeta | undefined;
	onPageChange: (page: number) => void;
}) {
	return (
		<section className="flex flex-col gap-4">
			<DataTable
				columns={getMemberColumns(canManage)}
				data={members}
				toolbar={() => <MembersToolbar />}
				isLoading={isPending}
				errorMessage={
					isError ? (errorMessage ?? "Failed to load members.") : undefined
				}
				emptyMessage="No members match your filters."
			/>

			{meta && (
				<Pagination
					page={meta.page}
					totalPages={meta.total_pages}
					onPageChange={onPageChange}
				/>
			)}
		</section>
	);
}

function InvitesTable({
	invites,
	isPending,
	isError,
	errorMessage,
}: {
	invites: Invite[];
	isPending: boolean;
	isError: boolean;
	errorMessage?: string;
}) {
	return (
		<DataTable
			columns={inviteColumns}
			data={invites}
			isLoading={isPending}
			errorMessage={
				isError ? (errorMessage ?? "Failed to load invites.") : undefined
			}
			emptyMessage="No pending invites"
		/>
	);
}
