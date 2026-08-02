import { MoreHorizontal } from "lucide-react";
import {
	Button,
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuLabel,
	DropdownMenuTrigger,
} from "#/components/ui";
import type { Invite } from "#/features/members/types";
import { canManageInvites } from "#/features/members/utils";
import { useModalStore } from "#/stores/useModalStore";
import { useWorkspaceStore } from "#/stores/useWorkspaceStore";

type InviteActionsProps = {
	invite: Invite;
};

// Row action menu for pending invites. Owner/admin only — matches
// RequireRole on DELETE /members/invites/:id.
export function InviteActions({ invite }: InviteActionsProps) {
	const myRole = useWorkspaceStore((s) => s.active?.role);
	const openModal = useModalStore((s) => s.openModal);

	if (!canManageInvites(myRole)) {
		return null;
	}

	const onModalOpen = () => {
		openModal({ type: "cancel-invite", props: { invite } });
	};
	return (
		<DropdownMenu>
			<DropdownMenuTrigger asChild>
				<Button variant="ghost" size="icon-sm" aria-label="Invite actions">
					<MoreHorizontal />
				</Button>
			</DropdownMenuTrigger>
			<DropdownMenuContent align="end" className="w-44">
				<DropdownMenuLabel>Actions</DropdownMenuLabel>
				<DropdownMenuItem variant="destructive" onSelect={onModalOpen}>
					Cancel invite
				</DropdownMenuItem>
			</DropdownMenuContent>
		</DropdownMenu>
	);
}
