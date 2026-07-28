import {
  InputOTP,
  InputOTPGroup,
  InputOTPSlot,
} from "#/components/ui/index.ts";
import { normalizeUserCode } from "../schemas/index.ts";

// Segmented XXXX-XXXX entry for the device code printed in the terminal. Eight
// mono slots split 4 + 4 by a hard dash. Controlled so it drops into a
// react-hook-form <Controller/>; every change is normalised (uppercased, purged
// of characters the alphabet never uses), so pasting "wdjf-x4k2" just works.
function DeviceCodeInput({
  value,
  onChange,
  onComplete,
  "aria-invalid": ariaInvalid,
}: {
  value: string;
  onChange: (value: string) => void;
  onComplete?: (value: string) => void;
  "aria-invalid"?: boolean;
}) {
  const slotClass =
    "h-14 flex-1 rounded-md! border-l! font-mono text-xl font-semibold uppercase";

  return (
    <InputOTP
      maxLength={8}
      value={value}
      onChange={(next) => onChange(normalizeUserCode(next))}
      onComplete={onComplete}
      containerClassName="justify-center gap-3"
      // Screen-reader label — the visual heading sits above, not on the input.
      aria-label="Device code"
    >
      <InputOTPGroup className="w-full gap-2">
        {[0, 1, 2, 3].map((i) => (
          <InputOTPSlot
            key={i}
            index={i}
            aria-invalid={ariaInvalid}
            className={slotClass}
          />
        ))}
      </InputOTPGroup>

      {/* Hard separator that echoes the printed WDJF-X4K2 form. */}
      <span
        aria-hidden="true"
        className="h-px w-3 shrink-0 rounded-full bg-muted-foreground/50"
      />

      <InputOTPGroup className="w-full gap-2">
        {[4, 5, 6, 7].map((i) => (
          <InputOTPSlot
            key={i}
            index={i}
            aria-invalid={ariaInvalid}
            className={slotClass}
          />
        ))}
      </InputOTPGroup>
    </InputOTP>
  );
}

export { DeviceCodeInput };
