import { describe, it, expect } from "vitest";
import { gameState } from "./gameState";
import { v4 } from "uuid";

describe("Game state tests", () => {
  it("Should be able to setup a game state", () => {
    expect(gameState).toBeTruthy();
    expect(gameState.validate()).toBeFalsy();

    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid);
    expect(gameState.validate()).toBeTruthy();
  });
});
