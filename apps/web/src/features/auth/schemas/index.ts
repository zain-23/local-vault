import { z } from "zod";

// Validation schemas for the auth forms. Each form infers its value type from
// the matching schema (single source of truth), and react-hook-form validates
// against these via @hookform/resolvers/zod.

const email = z.email({ message: "Enter a valid email address." });

export const signupSchema = z.object({
  name: z.string().min(1, "Enter your name."),
  email,
  password: z.string().min(12, "Use at least 12 characters."),
  // Must be checked — an unticked box fails validation (stays a boolean type).
  terms: z.boolean().refine((v) => v, {
    message: "Accept the Terms and Privacy Policy to continue.",
  }),
});

export const loginSchema = z.object({
  email,
  password: z.string().min(1, "Enter your password."),
});

export const forgotPasswordSchema = z.object({ email });

export const magicLinkSchema = z.object({ email });

export const resetPasswordSchema = z
  .object({
    password: z.string().min(12, "Use at least 12 characters."),
    confirmPassword: z.string(),
  })
  // Cross-field check: surface the mismatch on the confirm field.
  .refine((data) => data.password === data.confirmPassword, {
    message: "Passwords do not match.",
    path: ["confirmPassword"],
  });

export const twoFactorSchema = z.object({
  code: z.string().regex(/^\d{6}$/, "Enter the 6-digit code."),
});

export type LoginValues = z.infer<typeof loginSchema>;
export type SignupValues = z.infer<typeof signupSchema>;
export type ForgotPasswordValues = z.infer<typeof forgotPasswordSchema>;
export type MagicLinkValues = z.infer<typeof magicLinkSchema>;
export type ResetPasswordValues = z.infer<typeof resetPasswordSchema>;
export type TwoFactorValues = z.infer<typeof twoFactorSchema>;
