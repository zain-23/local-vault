import {
	Button,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "#/components/ui";
import { useRemoveMember } from "#/features/members/hooks";
import type { Member } from "#/features/members/types";
import { useModalStore } from "#/stores/useModalStore";

// Destructive confirm — DELETE /workspaces/:wid/members/:userId.
export function RemoveMemberModal() {
	const { props, closeModal } = useModalStore();
	const member = props.member as Member | undefined;
	const remove = useRemoveMember();

	if (!member) return null;

	const userId = member.user_id;

	function handleRemove() {
		remove.mutate(userId, {
			onSuccess: () => closeModal(),
		});
	}

	return (
		<DialogContent>
			<DialogHeader>
				<DialogTitle>Remove member</DialogTitle>
				<DialogDescription>
					Remove{" "}
					<span className="font-medium text-foreground">{member.name}</span>{" "}
					from this workspace? They'll lose access immediately. This can't be
					undone.
				</DialogDescription>
			</DialogHeader>

			<DialogFooter>
				<Button
					variant="outline"
					onClick={closeModal}
					disabled={remove.isPending}
				>
					Cancel
				</Button>
				<Button
					variant="destructive"
					onClick={handleRemove}
					isLoading={remove.isPending}
				>
					Remove member
				</Button>
			</DialogFooter>
		</DialogContent>
	);
}
