import type { ComponentType } from "react";
import { Dialog } from "#/components/ui";
import { ChangeRoleModal } from "#/features/members/modals/ChangeRoleModal.tsx";
import { InviteMemberModal } from "#/features/members/modals/InviteMemberModal.tsx";
import { RemoveMemberModal } from "#/features/members/modals/RemoveMemberModal.tsx";
import { type ModalType, useModalStore } from "#/stores/useModalStore";

// Registry: maps a modal type to the component that renders its DialogContent.
// Register a new modal by adding a ModalType and an entry here — feature code
// never imports Dialog; it just calls openModal({ type }).
const MODAL_REGISTRY: Record<ModalType, ComponentType> = {
	"invite-member": InviteMemberModal,
	"change-role": ChangeRoleModal,
	"remove-member": RemoveMemberModal,
};

// Mounted once in the protected shell. Owns the Dialog root + open state; the
// active registry component supplies the content.
export function GlobalModalProvider() {
	const { type, closeModal } = useModalStore();
	const ActiveModal = type ? MODAL_REGISTRY[type] : null;

	return (
		<Dialog open={type !== null} onOpenChange={(open) => !open && closeModal()}>
			{ActiveModal ? <ActiveModal /> : null}
		</Dialog>
	);
}
