import { useState } from "react";
import { toast } from "sonner";
import {
	Button,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	Field,
	FieldLabel,
	Select,
	SelectContent,
	SelectItem,
	SelectTrigger,
	SelectValue,
} from "#/components/ui";
import { ASSIGNABLE_ROLES, type Member } from "#/features/members/types";
import { useModalStore } from "#/stores/useModalStore";

// Change-role modal — maps to PUT /workspaces/:wid/members/:userId/role. Only
// admin/member are assignable (the server rejects promoting to owner).
export function ChangeRoleModal() {
	const { props, closeModal } = useModalStore();
	const member = props.member as Member | undefined;
	const initial = member?.role === "admin" ? "admin" : "member";
	const [role, setRole] = useState<(typeof ASSIGNABLE_ROLES)[number]>(initial);

	if (!member) return null;

	function handleSubmit(e: React.FormEvent) {
		e.preventDefault();
		// TODO: wire to the real endpoint. Mock for now.
		toast.success(`${member?.name} is now ${role}`);
		closeModal();
	}

	return (
		<DialogContent>
			<DialogHeader>
				<DialogTitle>Change role</DialogTitle>
				<DialogDescription>
					Update the workspace role for {member.name}.
				</DialogDescription>
			</DialogHeader>

			<form
				id="change-role-form"
				onSubmit={handleSubmit}
				className="grid gap-4"
			>
				<Field>
					<FieldLabel htmlFor="change-role-select">Role</FieldLabel>
					<Select
						value={role}
						onValueChange={(v) =>
							setRole(v as (typeof ASSIGNABLE_ROLES)[number])
						}
					>
						<SelectTrigger id="change-role-select" className="w-full">
							<SelectValue />
						</SelectTrigger>
						<SelectContent>
							{ASSIGNABLE_ROLES.map((r) => (
								<SelectItem key={r} value={r} className="capitalize">
									{r}
								</SelectItem>
							))}
						</SelectContent>
					</Select>
				</Field>
			</form>

			<DialogFooter>
				<Button variant="outline" onClick={closeModal}>
					Cancel
				</Button>
				<Button
					type="submit"
					form="change-role-form"
					disabled={role === member.role}
				>
					Save changes
				</Button>
			</DialogFooter>
		</DialogContent>
	);
}
