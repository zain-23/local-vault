import { useQueryStates } from "nuqs";
import { ProfileSection } from "#/features/settings/components/ProfileSection.tsx";
import { SecuritySection } from "#/features/settings/components/SecuritySection.tsx";
import { SessionsSection } from "#/features/settings/components/SessionsSection.tsx";
import { SettingsNav } from "#/features/settings/components/SettingsNav.tsx";
import {
  type SettingsSection,
  settingsSearchOptions,
  settingsSearchParams,
} from "#/features/settings/utils";

export function SettingsPage() {
  const [{ section }, setParams] = useQueryStates(
    settingsSearchParams,
    settingsSearchOptions,
  );

  return (
    <div className="flex flex-col gap-6 p-6 max-w-7xl mx-auto">
      <div>
        <h1 className="text-xl font-semibold tracking-tight">Settings</h1>
        <p className="mt-1 text-sm text-muted-foreground">
          Manage your account
        </p>
      </div>

      <div className="flex flex-col gap-6 sm:flex-row sm:items-start sm:gap-10">
        <SettingsNav
          value={section}
          onChange={(next) => {
            void setParams({ section: next });
          }}
        />

        <div className="min-w-0 flex-1">
          <SettingsPanel section={section} />
        </div>
      </div>
    </div>
  );
}

function SettingsPanel({ section }: { section: SettingsSection }) {
  switch (section) {
    case "security":
      return <SecuritySection />;
    case "sessions":
      return <SessionsSection />;
    default:
      return <ProfileSection />;
  }
}
