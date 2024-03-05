// @vitest-environment node
import WebSocket from "isomorphic-ws";
import { v4 } from "uuid";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { wsBaseUrl } from "../apiClient";
import { gameState } from "./gameState";

describe("Game state tests", () => {
  vi.mock("isomorphic-ws", () => {
    const wsConstructor = vi.fn(function con() {
      this.send = vi.fn();
      this.onclose = vi.fn();
      this.onopen = vi.fn();
      this.onmessage = vi.fn();
      this.onerror = vi.fn();
    });

    return { default: wsConstructor };
  });

  beforeEach(() => {
    vi.resetAllMocks();
  });

  it("Should be able to setup a game state", () => {
    expect(gameState).toBeTruthy();
    expect(gameState.validate()).toBeFalsy();

    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid);
    expect(gameState.validate()).toBeTruthy();
    expect(WebSocket).toHaveBeenCalledWith(
      `${wsBaseUrl}?game_id=${gid}&player_id=${pid}`,
    );
  });
});
