import { createIsomorphicFn } from "@tanstack/react-start";
import { getRequestHeader } from "@tanstack/react-start/server";
import { ApiError } from "./errors.ts";
import type { ApiResponse, HttpMethod, RequestOptions } from "./types.ts";

// Host is overridable per environment; the /api/v1 prefix is fixed to the server.
const BASE_URL = import.meta.env.VITE_API_URL as string | undefined;

// During SSR we forward the incoming request's cookies; in the browser the
// cookie jar attaches them for us, so there is nothing to read.
const getCookieHeader = createIsomorphicFn()
  .client(() => undefined)
  .server(() => getRequestHeader("cookie"));

// Refresh is browser-only: refresh_token is scoped to Path=/api/v1/auth, so the
// browser never sends it to the SSR server and there is nothing to forward there.
const isServer = createIsomorphicFn()
  .client(() => false)
  .server(() => true);

const REFRESH_PATH = "/auth/refresh";

export interface ApiClient {
  get<T>(path: string, options?: RequestOptions): Promise<ApiResponse<T>>;
  post<T>(
    path: string,
    body?: unknown,
    options?: RequestOptions,
  ): Promise<ApiResponse<T>>;
  patch<T>(
    path: string,
    body?: unknown,
    options?: RequestOptions,
  ): Promise<ApiResponse<T>>;
  put<T>(
    path: string,
    body?: unknown,
    options?: RequestOptions,
  ): Promise<ApiResponse<T>>;
  delete<T>(path: string, options?: RequestOptions): Promise<ApiResponse<T>>;
}

export class HttpClient implements ApiClient {
  // Held while a refresh is in flight so N concurrent 401s cause exactly one call.
  private refreshInFlight: Promise<boolean> | null = null;

  // baseUrl is injectable so tests can point at a stub server.
  constructor(private readonly baseUrl: string | undefined = BASE_URL) {}

  get<T>(path: string, options?: RequestOptions) {
    return this.request<T>("GET", path, undefined, options);
  }

  post<T>(path: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>("POST", path, body, options);
  }

  patch<T>(path: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>("PATCH", path, body, options);
  }

  put<T>(path: string, body?: unknown, options?: RequestOptions) {
    return this.request<T>("PUT", path, body, options);
  }

  delete<T>(path: string, options?: RequestOptions) {
    return this.request<T>("DELETE", path, undefined, options);
  }

  // Builds an absolute URL and appends defined query params.
  private buildUrl(path: string, params?: RequestOptions["params"]): string {
    if (!this.baseUrl) {
      throw new Error("Api key is not configured");
    }
    const url = new URL(`${this.baseUrl}${path}`);
    if (params) {
      for (const [key, value] of Object.entries(params)) {
        if (value !== undefined) url.searchParams.set(key, String(value));
      }
    }
    return url.toString();
  }

  // Asks the server for a new access token. Resolves true when the cookie was renewed.
  private refreshSession(): Promise<boolean> {
    this.refreshInFlight ??= (async () => {
      try {
        const res = await fetch(this.buildUrl(REFRESH_PATH), {
          method: "POST",
          credentials: "include",
          headers: { Accept: "application/json" },
        });
        // The response body is empty — the new access_token cookie is the payload.
        return res.ok;
      } catch {
        return false; // network failure — no token to retry with
      } finally {
        this.refreshInFlight = null; // let the next 401 start a fresh attempt
      }
    })();
    return this.refreshInFlight;
  }

  private async request<T>(
    method: HttpMethod,
    path: string,
    body?: unknown,
    options: RequestOptions = {},
    canRefresh = true, // false on the replay, so one failure can't loop
  ): Promise<ApiResponse<T>> {
    const { params, headers, ...init } = options;
    const cookie = getCookieHeader();

    const res = await fetch(this.buildUrl(path, params), {
      method,
      credentials: "include",
      headers: {
        Accept: "application/json",
        ...(body !== undefined ? { "Content-Type": "application/json" } : {}),
        ...(cookie ? { Cookie: cookie } : {}),
        ...headers,
      },
      body: body !== undefined ? JSON.stringify(body) : undefined,
      ...init,
    });

    // Expired access token — renew once and replay. Skipped on the server (no
    // refresh_token to send) and on the refresh call itself (would recurse).
    if (
      res.status === 401 &&
      canRefresh &&
      !isServer() &&
      path !== REFRESH_PATH &&
      (await this.refreshSession())
    ) {
      return this.request<T>(method, path, body, options, false);
    }

    // No content — synthesize an empty success envelope.
    if (res.status === 204) {
      return { data: undefined as T, success: true, message: "", status: 204 };
    }

    let raw: ApiResponse<T> & { error?: string };
    try {
      raw = (await res.json()) as ApiResponse<T> & { error?: string };
    } catch {
      throw new ApiError(res.ok ? 500 : res.status, "Internal server error");
    }

    // Errors come back as { error }. Throw so callers can branch on status, and so
    // react-query's error path means what it says.
    if (raw.error) throw new ApiError(res.status, raw.error);
    if (!res.ok)
      throw new ApiError(res.status, raw.message || "Request failed");

    return raw;
  }
}

// Shared singleton for app use; construct with a custom baseUrl in tests.
export const api = new HttpClient();
