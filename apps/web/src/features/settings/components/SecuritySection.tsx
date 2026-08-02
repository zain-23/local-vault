import { useQuery } from "@tanstack/react-query";
import { useState } from "react";
import QRCode from "react-qr-code";
import toast from "react-hot-toast";

import { Badge, Button, Input, Label } from "#/components/ui";
import { meQuery } from "#/features/auth/api";
import { SettingsBlock } from "#/features/settings/components/SettingsBlock.tsx";
import {
	useChangePassword,
	useDisable2FA,
	useEnable2FA,
	useRevokeOtherSessions,
	useVerify2FA,
} from "#/features/settings/hooks";
import { cn } from "#/lib/utils.ts";

export function SecuritySection() {
	const { data: user } = useQuery(meQuery);

	const [totpCode, setTotpCode] = useState("");
	const [disableCode, setDisableCode] = useState("");
	const [currentPassword, setCurrentPassword] = useState("");
	const [newPassword, setNewPassword] = useState("");
	const [confirmPassword, setConfirmPassword] = useState("");
	const [backupCodes, setBackupCodes] = useState<string[] | null>(null);

	const changePassword = useChangePassword();
	const enable2FA = useEnable2FA();
	const verify2FA = useVerify2FA();
	const disable2FA = useDisable2FA();
	const revokeOtherSessions = useRevokeOtherSessions();

	const twoFactorEnabled = user?.two_factor_enabled ?? false;
	const setup = enable2FA.data?.data;

	function updatePassword() {
		if (!currentPassword || !newPassword) {
			toast.error("Fill in current and new password");
			return;
		}
		if (newPassword !== confirmPassword) {
			toast.error("New passwords do not match");
			return;
		}
		if (newPassword.length < 8) {
			toast.error("Password must be at least 8 characters");
			return;
		}
		changePassword.mutate(
			{ current_password: currentPassword, new_password: newPassword },
			{
				onSuccess: () => {
					setCurrentPassword("");
					setNewPassword("");
					setConfirmPassword("");
				},
			},
		);
	}

	function confirm2FA() {
		if (totpCode.trim().length !== 6) {
			toast.error("Enter the 6-digit code from your app");
			return;
		}
		verify2FA.mutate(totpCode.trim(), {
			onSuccess: (response) => {
				setBackupCodes(response.data.backup_codes);
				setTotpCode("");
				enable2FA.reset();
				toast.success("Two-factor authentication enabled");
			},
		});
	}

	function confirmDisable2FA() {
		if (!disableCode.trim()) {
			toast.error("Enter a code from your authenticator app or backup codes");
			return;
		}
		const code = disableCode.trim();
		disable2FA.mutate(
			/^\d{6}$/.test(code) ? { totp_code: code } : { backup_code: code },
			{
				onSuccess: () => {
					setDisableCode("");
					setBackupCodes(null);
				},
			},
		);
	}

	return (
		<div>
			<div className="mb-2">
				<h2 className="text-[15px] font-semibold tracking-tight">Security</h2>
				<p className="mt-1 text-sm text-muted-foreground">
					Password and two-factor authentication
				</p>
			</div>

			<SettingsBlock
				label="Password"
				description="Use at least 8 characters. Updating signs out other sessions."
			>
				<Field
					id="current-password"
					label="Current password"
					type="password"
					autoComplete="current-password"
					placeholder="••••••••"
					value={currentPassword}
					onChange={setCurrentPassword}
				/>
				<Field
					id="new-password"
					label="New password"
					type="password"
					autoComplete="new-password"
					placeholder="At least 8 characters"
					value={newPassword}
					onChange={setNewPassword}
				/>
				<Field
					id="confirm-password"
					label="Confirm"
					type="password"
					autoComplete="new-password"
					placeholder="Repeat new password"
					value={confirmPassword}
					onChange={setConfirmPassword}
				/>
				<div className="flex justify-end">
					<Button
						isLoading={changePassword.isPending}
						disabled={changePassword.isPending}
						onClick={updatePassword}
					>
						Update password
					</Button>
				</div>
			</SettingsBlock>

			<SettingsBlock
				label="Two-factor authentication"
				description="Require a TOTP code from your authenticator app when signing in."
			>
				<div className="rounded-md border border-border bg-muted/40 px-3 py-2.5">
					<div className="flex flex-wrap items-center gap-2">
						<span className="text-sm font-medium">Authenticator app</span>
						<Badge
							variant="outline"
							className={cn(
								twoFactorEnabled
									? "border-success/40 text-success"
									: "text-muted-foreground",
							)}
						>
							{twoFactorEnabled ? "Enabled" : "Off"}
						</Badge>
					</div>
					<p className="mt-1 text-xs text-muted-foreground">
						{twoFactorEnabled
							? "Protects login with a one-time code"
							: "Not set up yet"}
					</p>
				</div>

				{twoFactorEnabled ? (
					<div className="flex flex-col gap-3">
						<Field
							id="totp-disable"
							label="Authenticator or backup code"
							placeholder="Enter a verification code"
							inputClassName="font-mono"
							value={disableCode}
							onChange={setDisableCode}
						/>
						<div className="flex justify-end">
							<Button
								size="sm"
								variant="destructive"
								disabled={disable2FA.isPending}
								onClick={confirmDisable2FA}
							>
								{disable2FA.isPending ? "Disabling..." : "Disable"}
							</Button>
						</div>
					</div>
				) : setup ? (
					<div className="flex flex-col gap-3 rounded-lg border border-border bg-muted/30 p-3.5">
						<p className="text-center text-xs text-muted-foreground">
							Scan this QR code with your authenticator app, then enter its
							code.
						</p>
						<div className="mx-auto rounded-lg bg-white p-3">
							<QRCode
								value={setup.otpauth_url}
								size={168}
								title="LocalVault authenticator setup"
							/>
						</div>
						<details className="text-center">
							<summary className="cursor-pointer text-xs text-muted-foreground hover:text-foreground">
								Can&apos;t scan? Enter the setup key instead
							</summary>
							<p className="mt-2 break-all font-mono text-xs text-foreground">
								{setup.secret}
							</p>
							<a
								className="mt-2 inline-block text-xs text-primary underline underline-offset-4"
								href={setup.otpauth_url}
							>
								Open in authenticator app
							</a>
						</details>
						<Field
							id="totp-setup"
							label="Verification code"
							placeholder="000000"
							inputClassName="font-mono tracking-widest"
							value={totpCode}
							onChange={setTotpCode}
						/>
						<Button disabled={verify2FA.isPending} onClick={confirm2FA}>
							{verify2FA.isPending ? "Verifying..." : "Confirm and enable"}
						</Button>
						<Button
							size="sm"
							variant="outline"
							onClick={() => enable2FA.reset()}
						>
							Cancel setup
						</Button>
					</div>
				) : (
					<Button
						isLoading={enable2FA.isPending}
						onClick={() => enable2FA.mutate()}
					>
						Enable 2FA
					</Button>
				)}

				{backupCodes ? (
					<div className="rounded-lg border border-warning/40 bg-warning/10 p-3.5">
						<p className="text-sm font-medium">Save these backup codes now</p>
						<p className="mt-1 text-xs text-muted-foreground">
							They are shown only once and each can be used once.
						</p>
						<div className="mt-3 grid grid-cols-2 gap-2 font-mono text-xs">
							{backupCodes.map((code) => (
								<span key={code}>{code}</span>
							))}
						</div>
					</div>
				) : null}
			</SettingsBlock>

			<SettingsBlock
				danger
				label="Sign out other sessions"
				description="Revoke every session except this one. You'll stay signed in here."
			>
				<Button
					variant="destructive"
					isLoading={revokeOtherSessions.isPending}
					onClick={() => revokeOtherSessions.mutate()}
				>
					Sign out other sessions
				</Button>
			</SettingsBlock>
		</div>
	);
}

function Field({
	id,
	label,
	type = "text",
	value,
	onChange,
	placeholder,
	autoComplete,
	inputClassName,
}: {
	id: string;
	label: string;
	type?: string;
	value: string;
	onChange: (value: string) => void;
	placeholder?: string;
	autoComplete?: string;
	inputClassName?: string;
}) {
	return (
		<div className="flex flex-col gap-1.5">
			<Label htmlFor={id} className="text-xs text-muted-foreground">
				{label}
			</Label>
			<Input
				id={id}
				type={type}
				value={value}
				placeholder={placeholder}
				autoComplete={autoComplete}
				className={inputClassName}
				onChange={(event) => onChange(event.target.value)}
			/>
		</div>
	);
}
