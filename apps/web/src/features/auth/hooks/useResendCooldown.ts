import { useEffect, useState } from "react";

import { RESEND_COOLDOWN_SECONDS } from "../utils/constant.ts";

/**
 * In-memory resend cooldown. Starts armed so the check-email screen can't
 * immediately fire another request after the form that landed here.
 */
export function useResendCooldown(durationSeconds = RESEND_COOLDOWN_SECONDS) {
	const [endsAt, setEndsAt] = useState(
		() => Date.now() + durationSeconds * 1000,
	);
	const [now, setNow] = useState(() => Date.now());

	useEffect(() => {
		const id = window.setInterval(() => setNow(Date.now()), 50);
		return () => window.clearInterval(id);
	}, []);

	const msLeft = Math.max(0, endsAt - now);
	const secondsLeft = Math.ceil(msLeft / 1000);
	const progress = Math.min(1, msLeft / (durationSeconds * 1000));
	const isCoolingDown = msLeft > 0;

	const restart = () => {
		const next = Date.now() + durationSeconds * 1000;
		setEndsAt(next);
		setNow(Date.now());
	};

	return { secondsLeft, progress, isCoolingDown, restart, durationSeconds };
}
