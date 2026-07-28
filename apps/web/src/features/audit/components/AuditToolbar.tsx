import { useQueryStates } from "nuqs";

import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "#/components/ui";
import {
  ACTION_PREFIX_LABELS,
  ACTION_PREFIXES,
  type ActionPrefixFilter,
  auditSearchOptions,
  auditSearchParams,
  RANGE_PRESET_LABELS,
  RANGE_PRESETS,
  type RangePreset,
} from "#/features/audit/utils";

export function AuditToolbar() {
  const [{ action_prefix, range }, setParams] = useQueryStates(
    auditSearchParams,
    auditSearchOptions,
  );

  return (
    <div className="flex flex-wrap items-center gap-2">
      <Select
        value={action_prefix}
        onValueChange={(next) =>
          setParams({
            action_prefix: next as ActionPrefixFilter,
            page: 1,
          })
        }
      >
        <SelectTrigger className="w-35" aria-label="Filter by action">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {ACTION_PREFIXES.map((prefix) => (
            <SelectItem key={prefix} value={prefix}>
              {ACTION_PREFIX_LABELS[prefix]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>

      <Select
        value={range}
        onValueChange={(next) =>
          setParams({
            range: next as RangePreset,
            page: 1,
          })
        }
      >
        <SelectTrigger className="w-37" aria-label="Filter by time range">
          <SelectValue />
        </SelectTrigger>
        <SelectContent>
          {RANGE_PRESETS.map((preset) => (
            <SelectItem key={preset} value={preset}>
              {RANGE_PRESET_LABELS[preset]}
            </SelectItem>
          ))}
        </SelectContent>
      </Select>
    </div>
  );
}
