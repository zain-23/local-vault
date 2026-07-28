// Thrown by the api client for any non-success response. Carries the HTTP status
// so callers can tell the difference between answers that look alike otherwise —
// a 401 means "not logged in", a 500 means "we don't know".
export class ApiError extends Error {
	constructor(
		readonly status: number,
		message: string,
	) {
		super(message);
		this.name = "ApiError";
	}
}
