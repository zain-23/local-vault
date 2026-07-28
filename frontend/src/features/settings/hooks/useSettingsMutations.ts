import { useMutation, useQueryClient } from "@tanstack/react-query";
import { toast } from "sonner";

import { AUTH_KEYS } from "#/features/auth/api";
import {
  SETTINGS_KEYS,
  settingsService,
  type ChangePasswordInput,
  type Disable2FAInput,
  type UpdateProfileInput,
} from "#/features/settings/api";

function invalidateAccount(queryClient: ReturnType<typeof useQueryClient>) {
  return queryClient.invalidateQueries({ queryKey: AUTH_KEYS.me() });
}

export function useUpdateProfile() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: SETTINGS_KEYS.updateProfile(),
    mutationFn: (input: UpdateProfileInput) =>
      settingsService.updateProfile(input),
    onSuccess: async (response) => {
      toast.success(response.message || "Profile updated");
      await invalidateAccount(queryClient);
    },
    onError: (error) =>
      toast.error(error.message || "Could not update profile"),
  });
}

export function useChangePassword() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: SETTINGS_KEYS.changePassword(),
    mutationFn: (input: ChangePasswordInput) =>
      settingsService.changePassword(input),
    onSuccess: async (response) => {
      toast.success(response.message || "Password updated");
      await queryClient.invalidateQueries({
        queryKey: SETTINGS_KEYS.sessions(),
      });
    },
    onError: (error) =>
      toast.error(error.message || "Could not update password"),
  });
}

export function useEnable2FA() {
  return useMutation({
    mutationKey: SETTINGS_KEYS.enable2FA(),
    mutationFn: () => settingsService.enable2FA(),
    onError: (error) =>
      toast.error(error.message || "Could not start 2FA setup"),
  });
}

export function useVerify2FA() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: SETTINGS_KEYS.verify2FA(),
    mutationFn: (totpCode: string) => settingsService.verify2FA(totpCode),
    onSuccess: async () => {
      await invalidateAccount(queryClient);
    },
    onError: (error) =>
      toast.error(error.message || "Could not verify the code"),
  });
}

export function useDisable2FA() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: SETTINGS_KEYS.disable2FA(),
    mutationFn: (input: Disable2FAInput) => settingsService.disable2FA(input),
    onSuccess: async (response) => {
      toast.success(response.message || "Two-factor authentication disabled");
      await invalidateAccount(queryClient);
    },
    onError: (error) => toast.error(error.message || "Could not disable 2FA"),
  });
}

export function useRevokeSession() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: SETTINGS_KEYS.revokeSession(),
    mutationFn: (sessionId: string) => settingsService.revokeSession(sessionId),
    onSuccess: async (response) => {
      toast.success(response.message || "Session revoked");
      await queryClient.invalidateQueries({
        queryKey: SETTINGS_KEYS.sessions(),
      });
    },
    onError: (error) =>
      toast.error(error.message || "Could not revoke session"),
  });
}

export function useRevokeOtherSessions() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationKey: SETTINGS_KEYS.revokeOtherSessions(),
    mutationFn: () => settingsService.revokeOtherSessions(),
    onSuccess: async (response) => {
      toast.success(response.message || "Other sessions signed out");
      await queryClient.invalidateQueries({
        queryKey: SETTINGS_KEYS.sessions(),
      });
    },
    onError: (error) =>
      toast.error(error.message || "Could not sign out sessions"),
  });
}
