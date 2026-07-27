import { queryOptions } from "@tanstack/react-query";

import { SETTINGS_KEYS } from "./settings.keys.ts";
import { settingsService } from "./settings.service.ts";

export const sessionsQuery = queryOptions({
	queryKey: SETTINGS_KEYS.sessions(),
	queryFn: async () => (await settingsService.listSessions()).data,
	staleTime: 30_000,
});
