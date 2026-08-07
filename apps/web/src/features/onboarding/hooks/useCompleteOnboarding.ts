import { useMutation, useQueryClient } from "@tanstack/react-query";
import { useNavigate } from "@tanstack/react-router";
import toast from "react-hot-toast";

import { meQuery } from "#/features/auth/api";
import { ONBOARDING_KEYS, onboardingService } from "#/features/onboarding/api";
import type { ApiResponse } from "#/services/api";

// Final step — flip the account's `onboarded` flag, then land the user in the
// app. We refetch meQuery (not just invalidate) and await it: the /onboarding
// guard reads that cached user, so a stale `onboarded: false` would bounce the
// user straight back here.
export function useCompleteOnboarding() {
	const navigate = useNavigate();
	const queryClient = useQueryClient();

	return useMutation<ApiResponse<null>, Error, void>({
		mutationFn: () => onboardingService.completeOnboarding(),
		mutationKey: ONBOARDING_KEYS.complete(),
		onSuccess: async () => {
			await queryClient.fetchQuery({ ...meQuery, staleTime: 0 });
			navigate({ to: "/dashboard" });
		},
		onError: (error) => {
			toast.error(error.message);
		},
	});
}
