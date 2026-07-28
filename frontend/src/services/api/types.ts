// The success envelope every endpoint returns — mirrors the Go server's
// `common/response.Envelope`. Error responses ({ error }) are normalized into
// this same shape by the client (success:false, message = the error text).
export interface ApiResponse<T> {
	data: T;
	success: boolean;
	message: string;
	status: number;
}

export type HttpMethod = "GET" | "POST" | "PATCH" | "PUT" | "DELETE";

// Per-request options: query params plus anything fetch accepts (signal,
// custom headers, cache…) except the parts the client controls.
export interface RequestOptions extends Omit<RequestInit, "method" | "body"> {
	params?: Record<string, string | number | boolean | undefined>;
}
