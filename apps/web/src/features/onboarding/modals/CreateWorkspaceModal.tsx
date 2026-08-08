import { zodResolver } from "@hookform/resolvers/zod";
import { Controller, useForm } from "react-hook-form";

import {
	Button,
	DialogContent,
	DialogDescription,
	DialogFooter,
	DialogHeader,
	DialogTitle,
	Field,
	FieldError,
	FieldLabel,
	Input,
} from "#/components/ui";
import { cn } from "#/lib/utils.ts";
import { useModalStore } from "#/stores/useModalStore";
import { useSaveWorkspace } from "../hooks/index.ts";
import { DEFAULT_WORKSPACE_ICON, WORKSPACE_ICONS, WorkspaceIcon } from "../lib/workspaceIcons.tsx";
import { type WorkspaceValues, workspaceSchema } from "../schemas/index.ts";

// Sidebar switcher's "create another workspace" action. Same shape as
// StepWorkspace, but always creates (no `existing` to prefill/rename) since
// this fires after onboarding, when the user already has at least one.
export function CreateWorkspaceModal() {
	const closeModal = useModalStore((s) => s.closeModal);

	const {
		register,
		control,
		handleSubmit,
		formState: { errors },
	} = useForm<WorkspaceValues>({
		resolver: zodResolver(workspaceSchema),
		defaultValues: { name: "", icon: DEFAULT_WORKSPACE_ICON },
	});

	const { isPending, mutate } = useSaveWorkspace(undefined, {
		onSuccess: closeModal,
	});

	const onSubmit = handleSubmit((values) => {
		mutate(values);
	});

	return (
		<DialogContent>
			<DialogHeader>
				<DialogTitle>Create workspace</DialogTitle>
				<DialogDescription>
					A new workspace has its own vaults, members, and audit history.
				</DialogDescription>
			</DialogHeader>

			<form
				id="create-workspace-form"
				onSubmit={onSubmit}
				noValidate
				className="grid gap-4"
			>
				<Field data-invalid={!!errors.name}>
					<FieldLabel htmlFor="create-workspace-name">
						Workspace name
					</FieldLabel>
					<Input
						id="create-workspace-name"
						placeholder="Acme Inc"
						autoFocus
						autoComplete="organization"
						aria-invalid={!!errors.name}
						{...register("name")}
					/>
					<FieldError errors={[errors.name]} />
				</Field>

				<Field data-invalid={!!errors.icon}>
					<FieldLabel id="create-workspace-icon-label">
						Choose a mark
					</FieldLabel>
					<Controller
						name="icon"
						control={control}
						render={({ field }) => (
							<div
								className="mt-2 grid grid-cols-6 gap-2"
								role="listbox"
								aria-labelledby="create-workspace-icon-label"
							>
								{WORKSPACE_ICONS.map((preset) => {
									const selected = field.value === preset.id;
									return (
										<button
											key={preset.id}
											type="button"
											role="option"
											aria-selected={selected}
											aria-label={preset.label}
											onClick={() => field.onChange(preset.id)}
											className={cn(
												"aspect-square w-full rounded-[10px] p-0 transition-[transform,box-shadow] duration-150",
												"focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 focus-visible:ring-offset-card",
												"hover:scale-[1.04] active:scale-[0.97]",
												selected
													? "ring-2 ring-primary ring-offset-2 ring-offset-card"
													: "ring-1 ring-transparent hover:ring-border",
											)}
										>
											<WorkspaceIcon
												id={preset.id}
												className="size-full rounded-[9px]"
											/>
										</button>
									);
								})}
							</div>
						)}
					/>
					<FieldError errors={[errors.icon]} />
				</Field>
			</form>

			<DialogFooter>
				<Button variant="outline" onClick={closeModal} disabled={isPending}>
					Cancel
				</Button>
				<Button
					type="submit"
					form="create-workspace-form"
					isLoading={isPending}
				>
					Create workspace
				</Button>
			</DialogFooter>
		</DialogContent>
	);
}
