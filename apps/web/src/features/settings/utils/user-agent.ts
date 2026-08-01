export interface ParsedUserAgent {
	browser: string;
	os: string;
}

function parseBrowser(ua: string): string {
	if (/Edg\//i.test(ua)) return "Edge";
	if (/OPR\/|Opera/i.test(ua)) return "Opera";
	if (/Firefox\//i.test(ua)) return "Firefox";
	if (/Chrome\//i.test(ua) && !/Chromium/i.test(ua)) return "Chrome";
	if (/Chromium\//i.test(ua)) return "Chromium";
	if (/Safari\//i.test(ua) && !/Chrome\//i.test(ua)) return "Safari";
	return "Unknown browser";
}

function parseOS(ua: string): string {
	if (/Android/i.test(ua)) return "Android";
	if (/iPhone|iPad|iPod/i.test(ua)) return "iOS";
	if (/Windows NT/i.test(ua)) return "Windows";
	if (/Mac OS X|Macintosh/i.test(ua)) return "macOS";
	if (/CrOS/i.test(ua)) return "Chrome OS";
	if (/Linux/i.test(ua)) return "Linux";
	return "Unknown OS";
}

export function parseUserAgent(userAgent: string): ParsedUserAgent {
	const ua = userAgent.trim();
	if (!ua) {
		return { browser: "Unknown browser", os: "Unknown OS" };
	}
	return {
		browser: parseBrowser(ua),
		os: parseOS(ua),
	};
}

export function formatSessionDevice(userAgent: string): string {
	const { browser, os } = parseUserAgent(userAgent);
	return `${browser} on ${os}`;
}
