import { useMutation, useQueryClient } from "@tanstack/react-query";
import toast from "react-hot-toast";

import { AUTH_KEYS } from "#/features/auth/api";
import {
	SETTINGS_KEYS,
	settingsService,
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
