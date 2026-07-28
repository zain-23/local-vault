import type { RangePreset } from "./search-params.ts";

/** Floor to the minute so `from` stays stable across renders within the same minute. */
function floorToMinute(d: Date): Date {
  const t = d.getTime();
  return new Date(t - (t % 60_000));
}

/** RFC3339 lower bound for a range preset, or undefined for unbounded. */
export function rangePresetToFrom(
  preset: RangePreset,
  now = new Date(),
): string | undefined {
  if (preset === "all") return undefined;
  const ms =
    preset === "24h"
      ? 24 * 60 * 60 * 1000
      : preset === "7d"
        ? 7 * 24 * 60 * 60 * 1000
        : 30 * 24 * 60 * 60 * 1000;
  return new Date(floorToMinute(now).getTime() - ms).toISOString();
}

export const RANGE_PRESET_LABELS: Record<RangePreset, string> = {
  "24h": "Last 24 hours",
  "7d": "Last 7 days",
  "30d": "Last 30 days",
  all: "All time",
};

export const ACTION_PREFIX_LABELS: Record<string, string> = {
  all: "Any action",
  vault: "Vault",
  member: "Member",
  workspace: "Workspace",
  device: "Device",
};
