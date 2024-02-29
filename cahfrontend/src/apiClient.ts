import { Api, HttpClient } from "./api";

const defaultBaseUrl = new HttpClient().baseUrl;

/**
 * Default API client which loads the base URL from the environment variables.
 */
export const apiClient = new Api({
  baseUrl: import.meta.env.VITE_API_BASE_URL ?? defaultBaseUrl,
});

/**
 * The URL for the (http upgraded) websocket for the RPC channel that the real
 * time game service uses.
 */
export const wsBaseUrl =
  import.meta.env.VITE_WS_BASE_URL ??
  "ws://" + defaultBaseUrl.split("://")[1] + "/games/join";
