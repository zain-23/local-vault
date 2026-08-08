import type { LucideIcon } from "lucide-react";
import {
	Boxes,
	Cloud,
	Database,
	Fingerprint,
	FolderLock,
	KeyRound,
	Lock,
	Rocket,
	Shield,
	Terminal,
	Vault,
	Wrench,
} from "lucide-react";

import { cn } from "#/lib/utils.ts";

export const DEFAULT_WORKSPACE_ICON = "vault";

export type WorkspaceIconId =
	| "vault"
	| "lock"
	| "key"
	| "shield"
	| "folder"
	| "rocket"
	| "wrench"
	| "database"
	| "cloud"
	| "terminal"
	| "boxes"
	| "fingerprint";

export type WorkspaceIconPreset = {
	id: WorkspaceIconId;
	label: string;
	icon: LucideIcon;
	bg: string;
};

export const WORKSPACE_ICONS: WorkspaceIconPreset[] = [
	{ id: "vault", label: "Vault", icon: Vault, bg: "#1a6b5c" },
	{ id: "lock", label: "Lock", icon: Lock, bg: "#c45c26" },
	{ id: "key", label: "Key", icon: KeyRound, bg: "#c9a227" },
	{ id: "shield", label: "Shield", icon: Shield, bg: "#2f5d8a" },
	{ id: "folder", label: "Folder", icon: FolderLock, bg: "#6b4f3a" },
	{ id: "rocket", label: "Rocket", icon: Rocket, bg: "#d97706" },
	{ id: "wrench", label: "Wrench", icon: Wrench, bg: "#4a5568" },
	{ id: "database", label: "Database", icon: Database, bg: "#0e7490" },
	{ id: "cloud", label: "Cloud", icon: Cloud, bg: "#2563eb" },
	{ id: "terminal", label: "Terminal", icon: Terminal, bg: "#166534" },
	{ id: "boxes", label: "Boxes", icon: Boxes, bg: "#9a3412" },
	{ id: "fingerprint", label: "Fingerprint", icon: Fingerprint, bg: "#7c2d12" },
];

const byId = new Map(WORKSPACE_ICONS.map((p) => [p.id, p]));

export function getWorkspaceIcon(id: string | undefined | null): WorkspaceIconPreset {
	if (id && byId.has(id as WorkspaceIconId)) {
		return byId.get(id as WorkspaceIconId)!;
	}
	return byId.get(DEFAULT_WORKSPACE_ICON)!;
}

export const WORKSPACE_ICON_IDS = WORKSPACE_ICONS.map((p) => p.id) as [
	WorkspaceIconId,
	...WorkspaceIconId[],
];

type WorkspaceIconProps = {
	id?: string | null;
	size?: "sm" | "md" | "lg";
	className?: string;
};

const sizeClass = {
	sm: "size-6 rounded-md",
	md: "size-8 rounded-md",
	lg: "size-11 rounded-[10px]",
} as const;

export function WorkspaceIcon({
	id,
	size = "md",
	className,
}: WorkspaceIconProps) {
	const preset = getWorkspaceIcon(id);
	const Icon = preset.icon;

	return (
		<div
			className={cn(
				"flex shrink-0 items-center justify-center text-white",
				sizeClass[size],
				className,
			)}
			style={{ backgroundColor: preset.bg }}
			aria-hidden
		>
			<Icon className="size-[45%]" strokeWidth={2.25} />
		</div>
	);
}
