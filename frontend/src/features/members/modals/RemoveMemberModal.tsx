import { toast } from "sonner";
import {
	Button,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
} from "#/components/ui";
import type { Member } from "#/features/members/types";
import { useModalStore } from "#/stores/useModalStore";

// Destructive confirm — maps to DELETE /workspaces/:wid/members/:userId.
export function RemoveMemberModal() {
	const { props, closeModal } = useModalStore();
	const member = props.member as Member | undefined;

	if (!member) return null;

	function handleRemove() {
		// TODO: wire to the real endpoint. Mock for now.
		toast.success(`${member?.name} was removed`);
		closeModal();
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
				<Button variant="outline" onClick={closeModal}>
					Cancel
				</Button>
				<Button variant="destructive" onClick={handleRemove}>
					Remove member
				</Button>
			</DialogFooter>
		</DialogContent>
	);
}
