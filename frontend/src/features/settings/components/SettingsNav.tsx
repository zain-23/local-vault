import { Button } from "#/components/ui";
import { SETTINGS_NAV, type SettingsSection } from "#/features/settings/utils";
import { cn } from "#/lib/utils.ts";

type SettingsNavProps = {
	value: SettingsSection;
	onChange: (section: SettingsSection) => void;
};

export function SettingsNav({ value, onChange }: SettingsNavProps) {
	return (
		<nav
			aria-label="Settings sections"
			className="flex w-full flex-col items-start gap-0.5 sm:w-40"
		>
			{SETTINGS_NAV.map((item) => {
				const active = value === item.id;
				return (
					<Button
						key={item.id}
						type="button"
						variant={active ? "secondary" : "ghost"}
						size="sm"
						onClick={() => onChange(item.id)}
						className={cn(
							"h-auto w-full justify-start p-1.5 text-left",
							active
								? "font-medium text-foreground"
								: "font-normal text-muted-foreground hover:text-foreground",
						)}
					>
						{item.label}
					</Button>
				);
			})}
		</nav>
	);
}
