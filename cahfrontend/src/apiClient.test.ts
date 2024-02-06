import { describe, it, expect } from "vitest";
import { wsBaseUrl } from "./apiClient";

describe("Api client", () => {
  it("Websocket should be defined", () => {
    expect(wsBaseUrl).toBeTruthy();
  });
});
