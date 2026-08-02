import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import toast from "react-hot-toast";

import { Button, Input, Label } from "#/components/ui";
import { meQuery } from "#/features/auth/api";
import { SettingsBlock } from "#/features/settings/components/SettingsBlock.tsx";
import { useUpdateProfile } from "#/features/settings/hooks";
import { formatMemberSince, initialsFromName } from "#/features/settings/utils";

export function ProfileSection() {
	const { data: user } = useQuery(meQuery);
	const updateProfile = useUpdateProfile();
	const [name, setName] = useState(user?.name ?? "");

	useEffect(() => {
		if (user?.name) setName(user.name);
	}, [user?.name]);

	const dirty = Boolean(user && name.trim() !== user.name);
	const initials = initialsFromName(name || user?.name || user?.email || "?");

	const onSave = () => {
		const trimmedName = name.trim();
		if (name.length < 2) {
			toast.error("Display name must be at least 2 characters");
			return;
		}
		updateProfile.mutate({ name: trimmedName });
	};

	return (
		<div>
			<SectionHeader title="Profile" subtitle="Your name and account email" />

			<SettingsBlock
				label="Avatar"
				description="Shown in the workspace and audit log."
			>
				<div className="flex items-center gap-3">
					<div
						aria-hidden
						className="flex size-10 shrink-0 items-center justify-center rounded-lg bg-primary/15 font-mono text-sm font-semibold text-primary"
					>
						{initials}
					</div>
					<span className="text-[13px] text-muted-foreground">
						Initials from your name
					</span>
				</div>
			</SettingsBlock>

			<SettingsBlock
				label="Display name"
				description="How teammates see you across LocalVault."
			>
				<Label htmlFor="settings-name" className="sr-only">
					Display name
				</Label>
				<Input
					id="settings-name"
					value={name}
					onChange={(e) => setName(e.target.value)}
					autoComplete="name"
				/>
			</SettingsBlock>

			<SettingsBlock
				label="Email"
				description="Used for login and vault invite codes. Contact support to change."
			>
				<Input value={user?.email ?? ""} disabled readOnly />
			</SettingsBlock>

			<SettingsBlock
				label="Member since"
				description="When this account was created."
			>
				<p className="font-mono text-sm text-muted-foreground">
					{user?.created_at ? formatMemberSince(user.created_at) : "—"}
				</p>
			</SettingsBlock>

			<div className="flex justify-end pt-4">
				<Button
					disabled={!dirty}
					isLoading={updateProfile.isPending}
					onClick={onSave}
				>
					Save changes
				</Button>
			</div>
		</div>
	);
}

function SectionHeader({
	title,
	subtitle,
}: {
	title: string;
	subtitle: string;
}) {
	return (
		<div className="mb-2">
			<h2 className="text-[15px] font-semibold tracking-tight">{title}</h2>
			<p className="mt-1 text-sm text-muted-foreground">{subtitle}</p>
		</div>
	);
}
