import type { ApiResponse, HttpMethod, RequestOptions } from "./types.ts";

// Host is overridable per environment; the /api/v1 prefix is fixed to the server.
const BASE_URL = import.meta.env.VITE_API_URL as string | undefined;

// Builds an absolute URL and appends defined query params.
function buildUrl(path: string, params?: RequestOptions["params"]): string {
  if (!BASE_URL) {
    throw new Error("Api key is not configured");
  }
  const url = new URL(`${BASE_URL}${path}`);
  if (params) {
    for (const [key, value] of Object.entries(params)) {
      if (value !== undefined) url.searchParams.set(key, String(value));
    }
  }
  return url.toString();
}

// Core request: serializes the body, sends auth cookies, parses the envelope,
// and returns the full ApiResponse<T> (data + message + status) or throws
// ApiError. The envelope is returned whole because several endpoints carry the
// meaningful text in `message` with `data` null (signup, forgot/reset, …).
async function request<T>(
  method: HttpMethod,
  path: string,
  body?: unknown,
  { params, headers, ...init }: RequestOptions = {},
): Promise<ApiResponse<T>> {
  const res = await fetch(buildUrl(path, params), {
    method,
    credentials: "include",
    headers: {
      Accept: "application/json",
      ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
      ...headers,
    },
    body: body !== undefined ? JSON.stringify(body) : undefined,
    ...init,
  });

  // No content — synthesize an empty success envelope.
  if (res.status === 204) {
    return { data: undefined as T, success: true, message: "", status: 204 };
  }

  let envelope: ApiResponse<T>;
  try {
    // Errors come back as { error }. Normalize that into the envelope: the error
    // text becomes `message` and success is false.
    const raw = (await res.json()) as ApiResponse<T> & { error?: string };
    envelope = raw.error
      ? {
          data: null as T,
          success: false,
          message: raw.error,
          status: res.status,
        }
      : raw;
  } catch {
    envelope = {
      data: null as T,
      message: "Internal server error",
      status: 500,
      success: false,
    };
  }

  return envelope;
}

// Reusable, fully type-safe HTTP client. Each method resolves to ApiResponse<T>.
// Usage:
//   const { data: user, message } = await api.post<User>("/auth/login", body);
//   const { data: vaults } = await api.get<Vault[]>("/vaults", { params });
export const api = {
  get: <T>(path: string, options?: RequestOptions) =>
    request<T>("GET", path, undefined, options),
  post: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>("POST", path, body, options),
  patch: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>("PATCH", path, body, options),
  put: <T>(path: string, body?: unknown, options?: RequestOptions) =>
    request<T>("PUT", path, body, options),
  delete: <T>(path: string, options?: RequestOptions) =>
    request<T>("DELETE", path, undefined, options),
};

// export type
export type ApiClient = typeof api;
