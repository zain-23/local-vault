import { parseAsInteger, parseAsStringLiteral } from "nuqs";

export const ACTION_PREFIXES = [
  "all",
  "vault",
  "member",
  "workspace",
  "device",
] as const;

export type ActionPrefixFilter = (typeof ACTION_PREFIXES)[number];

export const RANGE_PRESETS = ["24h", "7d", "30d", "all"] as const;

export type RangePreset = (typeof RANGE_PRESETS)[number];

export const auditSearchParams = {
  page: parseAsInteger.withDefault(1),
  action_prefix: parseAsStringLiteral(ACTION_PREFIXES).withDefault("all"),
  range: parseAsStringLiteral(RANGE_PRESETS).withDefault("7d"),
};

export const auditSearchOptions = {
  history: "replace" as const,
  shallow: true,
};
