import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery } from "@tanstack/react-query";
import { Controller, useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldLabel,
	Input,
} from "#/components/ui/index.ts";
import { cn } from "#/lib/utils.ts";
import { workspacesQuery } from "../api/index.ts";
import { useSaveWorkspace } from "../hooks/index.ts";
import {
	getWorkspaceIcon,
	WORKSPACE_ICONS,
	WorkspaceIcon,
} from "../lib/workspaceIcons.tsx";
import { type WorkspaceValues, workspaceSchema } from "../schemas/index.ts";
import { StepHeading } from "./OnboardingLayout.tsx";

function StepWorkspace({ onContinue }: { onContinue: () => void }) {
	// Prefetched by the /onboarding route loader, so this resolves synchronously
	// on first render — no empty-input flash. On a refresh after step 1 the
	// workspace already exists; we prefill from it and rename instead of creating
	// a duplicate. During onboarding the user has at most one workspace.
	const { data: workspaces } = useQuery(workspacesQuery);
	const existing = workspaces?.[0]?.workspace;

	const {
		register,
		control,
		handleSubmit,
		formState: { errors },
	} = useForm<WorkspaceValues>({
		resolver: zodResolver(workspaceSchema),
		// `values` (not `defaultValues`) so the field syncs once the query resolves.
		values: {
			name: existing?.name ?? "",
			icon: getWorkspaceIcon(existing?.icon).id,
		},
	});

	// The hook creates, renames, or no-ops based on `existing`; onContinue then
	// advances the wizard. Failures surface as a toast from the hook.
	const { isPending, mutate } = useSaveWorkspace(existing, {
		onSuccess: onContinue,
	});

	const onSubmit = handleSubmit((values) => {
		mutate(values);
	});

	return (
		<>
			<StepHeading
				title="Create your workspace"
				subtitle="Workspaces hold your team, vaults, and audit history."
			/>

			<form onSubmit={onSubmit} noValidate>
				{/* Identity strip — mark + name read as one unit, matching the sidebar switcher. */}
				<div className="flex items-end gap-3">
					<Field className="min-w-0 flex-1" data-invalid={!!errors.name}>
						<FieldLabel htmlFor="workspace-name">Workspace name</FieldLabel>
						<Input
							id="workspace-name"
							placeholder="Acme Inc"
							autoComplete="organization"
							aria-invalid={!!errors.name}
							{...register("name")}
						/>
						<FieldError errors={[errors.name]} />
					</Field>
				</div>

				<Field className="mt-6" data-invalid={!!errors.icon}>
					<FieldLabel id="workspace-icon-label">Choose a mark</FieldLabel>
					<Controller
						name="icon"
						control={control}
						render={({ field }) => (
							<div
								className="mt-2 grid grid-cols-6 gap-2"
								role="listbox"
								aria-labelledby="workspace-icon-label"
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

				<Button type="submit" className="mt-6 w-full" isLoading={isPending}>
					Continue
				</Button>
			</form>
		</>
	);
}

export { StepWorkspace };
