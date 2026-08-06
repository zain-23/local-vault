import { zodResolver } from "@hookform/resolvers/zod";
import { useQuery } from "@tanstack/react-query";
import { useForm } from "react-hook-form";

import {
	Button,
	Field,
	FieldError,
	FieldLabel,
	Input,
} from "#/components/ui/index.ts";
import { workspacesQuery } from "../api/index.ts";
import { useSaveWorkspace } from "../hooks/index.ts";
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
		handleSubmit,
		formState: { errors },
	} = useForm<WorkspaceValues>({
		resolver: zodResolver(workspaceSchema),
		// `values` (not `defaultValues`) so the field syncs once the query resolves.
		values: { name: existing?.name ?? "" },
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
				<Field data-invalid={!!errors.name}>
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

				<Button type="submit" className="mt-6 w-full" isLoading={isPending}>
					Continue
				</Button>
			</form>
		</>
	);
}

export { StepWorkspace };
