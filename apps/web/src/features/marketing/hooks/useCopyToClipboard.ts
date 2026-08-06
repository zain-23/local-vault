import { useCallback, useEffect, useRef, useState } from "react";

/**
 * Copies `text` and flips `copied` for `resetAfterMs`. Failures (insecure
 * origin, denied permission) leave `copied` false rather than lying to the user.
 */
export function useCopyToClipboard(text: string, resetAfterMs = 1500) {
	const [copied, setCopied] = useState(false);
	const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);

	// A copy right before unmount would otherwise setState on a dead component.
	useEffect(
		() => () => {
			if (timeoutRef.current) clearTimeout(timeoutRef.current);
		},
		[],
	);

	const copy = useCallback(async () => {
		try {
			await navigator.clipboard.writeText(text);
		} catch {
			return;
		}
		setCopied(true);
		if (timeoutRef.current) clearTimeout(timeoutRef.current);
		timeoutRef.current = setTimeout(() => setCopied(false), resetAfterMs);
	}, [text, resetAfterMs]);

	return { copied, copy };
}
