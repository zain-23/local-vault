import { useQuery } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toast } from "sonner";

import { Badge, Button, Input, Label } from "#/components/ui";
import { meQuery } from "#/features/auth/api";
import { SettingsBlock } from "#/features/settings/components/SettingsBlock.tsx";
import { cn } from "#/lib/utils.ts";

export function SecuritySection() {
  const { data: user } = useQuery(meQuery);
  const [twoFactorEnabled, setTwoFactorEnabled] = useState(false);
  const [setupOpen, setSetupOpen] = useState(false);
  const [totpCode, setTotpCode] = useState("");
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");

  useEffect(() => {
    if (user) setTwoFactorEnabled(user.two_factor_enabled);
  }, [user]);

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
            onClick={() => {
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
              setCurrentPassword("");
              setNewPassword("");
              setConfirmPassword("");
              toast.success("Password updated");
            }}
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
              ? "Protects login with a one-time code · 6 backup codes left"
              : "Not set up yet"}
          </p>
        </div>

        <div className="flex flex-wrap gap-2">
          {twoFactorEnabled ? (
            <>
              <Button
                size="sm"
                variant="outline"
                onClick={() => toast.message("Backup codes (UI only for now)")}
              >
                View backup codes
              </Button>
              <Button
                size="sm"
                variant="destructive"
                onClick={() => {
                  setTwoFactorEnabled(false);
                  setSetupOpen(false);
                  toast.success("Two-factor authentication disabled");
                }}
              >
                Disable
              </Button>
            </>
          ) : (
            <Button
              size="sm"
              variant={setupOpen ? "outline" : "default"}
              onClick={() => setSetupOpen((open) => !open)}
            >
              {setupOpen ? "Cancel setup" : "Enable 2FA"}
            </Button>
          )}
        </div>

        {setupOpen && !twoFactorEnabled ? (
          <div className="flex flex-col gap-3 rounded-lg border border-border bg-muted/30 p-3.5">
            <div className="mx-auto flex size-30 items-center justify-center rounded-lg border border-dashed border-border-strong font-mono text-[11px] text-muted-foreground">
              QR PLACEHOLDER
            </div>
            <p className="text-center text-xs text-muted-foreground">
              Secret{" "}
              <span className="font-mono text-foreground">
                JBSW Y3DP EHPK 3PXP
              </span>
            </p>
            <Field
              id="totp-setup"
              label="Verification code"
              placeholder="000000"
              inputClassName="font-mono tracking-widest"
              value={totpCode}
              onChange={setTotpCode}
            />
            <Button
              onClick={() => {
                if (totpCode.trim().length < 6) {
                  toast.error("Enter the 6-digit code from your app");
                  return;
                }
                setTwoFactorEnabled(true);
                setSetupOpen(false);
                setTotpCode("");
                toast.success("Two-factor authentication enabled");
              }}
            >
              Confirm and enable
            </Button>
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
          onClick={() => toast.success("Other sessions signed out")}
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
        onChange={(e) => onChange(e.target.value)}
      />
    </div>
  );
}
