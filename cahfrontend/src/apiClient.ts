import { CookieOptions } from "@solid-primitives/storage";

import { Api, HttpClient } from "./api";

// const defaultBaseUrl = new HttpClient().baseUrl;
const defaultBaseUrl = "http://localhost:3255/api";

/**
 * Default API client which loads the base URL from the environment variables.
 */
export const apiClient = new Api({
  baseUrl: import.meta.env.VITE_API_BASE_URL ?? defaultBaseUrl,
  customFetch: async (url, init) => {
    const resp = await fetch(url, {
      credentials: "include",
      mode: "cors",
      ...init,
    });
    return resp;
  },
});

/**
 * The URL for the (http upgraded) websocket for the RPC channel that the real
 * time game service uses.
 */
export const wsBaseUrl =
  import.meta.env.VITE_WS_BASE_URL ??
  "ws://localhost:3255/ws";

export const cookieOptions: Readonly<CookieOptions> = {
  path: "/",
  sameSite: "None",
  secure: true,
};
