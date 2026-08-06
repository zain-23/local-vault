import {
	Button,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "#/components/ui";
import { useCancelInvite } from "#/features/members/hooks";
import type { Invite } from "#/features/members/types";
import { useModalStore } from "#/stores/useModalStore";

// Destructive confirm — DELETE /workspaces/:wid/members/invites/:id.
export function CancelInviteModal() {
	const { props, closeModal } = useModalStore();
	const invite = props.invite as Invite | undefined;
	const cancel = useCancelInvite();

	if (!invite) return null;

	const inviteId = invite.id;

	function handleCancel() {
		cancel.mutate(inviteId, {
			onSuccess: () => closeModal(),
		});
	}

	return (
		<DialogContent>
			<DialogHeader>
				<DialogTitle>Cancel invite</DialogTitle>
				<DialogDescription>
					Cancel the invite for{" "}
					<span className="font-medium text-foreground">{invite.email}</span>?
					They won't be able to join with this link. This can't be undone.
				</DialogDescription>
			</DialogHeader>

			<DialogFooter>
				<Button
					variant="outline"
					onClick={closeModal}
					disabled={cancel.isPending}
				>
					Keep invite
				</Button>
				<Button
					variant="destructive"
					onClick={handleCancel}
					isLoading={cancel.isPending}
				>
					Cancel invite
				</Button>
			</DialogFooter>
		</DialogContent>
	);
}
