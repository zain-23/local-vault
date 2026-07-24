import { zodResolver } from "@hookform/resolvers/zod";
import { useNavigate } from "@tanstack/react-router";
import { Terminal } from "lucide-react";
import { Controller, useForm } from "react-hook-form";

import {
  Button,
  Field,
  FieldError,
  FieldGroup,
} from "#/components/ui/index.ts";
import {
  type DeviceCodeValues,
  deviceCodeSchema,
  toUserCodeParam,
} from "../schemas/index.ts";
import { DeviceCard } from "./DeviceCard.tsx";
import { DeviceCodeInput } from "./DeviceCodeInput.tsx";

// Step one of the browser flow (not in the original mockup): the user types the
// code their terminal printed, and we hand off to the approval screen. No
// network call happens here — the code is only checked for shape; the server
// validates it when the approval screen loads.
function SubmitCodeForm({ defaultCode = "" }: { defaultCode?: string }) {
  const navigate = useNavigate();
  const {
    control,
    handleSubmit,
    formState: { errors },
  } = useForm<DeviceCodeValues>({
    resolver: zodResolver(deviceCodeSchema),
    defaultValues: { code: defaultCode },
  });

  const onSubmit = handleSubmit(({ code }) =>
    navigate({
      to: "/device/confirmation",
      search: { user_code: toUserCodeParam(code) },
    }),
  );

  return (
    <DeviceCard
      icon={Terminal}
      title="Enter the code from your terminal"
      subtitle={
        <>
          LocalVault CLI showed you a code like{" "}
          <span className="font-mono text-foreground">WDJF-X4K2</span>. Type it
          here to link this device.
        </>
      }
    >
      <form onSubmit={onSubmit} noValidate>
        <FieldGroup className="gap-5">
          <Field data-invalid={!!errors.code}>
            <Controller
              control={control}
              name="code"
              render={({ field }) => (
                <DeviceCodeInput
                  value={field.value}
                  onChange={field.onChange}
                  // No auto-submit — the user reviews and clicks Continue.
                  aria-invalid={!!errors.code}
                />
              )}
            />
            <FieldError className="text-center" errors={[errors.code]} />
          </Field>

          <Button type="submit" className="w-full">
            Continue
          </Button>
        </FieldGroup>
      </form>

      <p className="mt-5 text-center text-xs text-muted-foreground">
        Didn't start this? You can safely close this page.
      </p>
    </DeviceCard>
  );
}

export { SubmitCodeForm };
