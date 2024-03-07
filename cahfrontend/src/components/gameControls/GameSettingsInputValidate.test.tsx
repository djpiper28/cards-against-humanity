import { Mock, describe, expect, it, vi } from "vitest";
import { v4 } from "uuid";
import {
  Settings,
  validateGamePassword,
  validateMaxGameRounds,
  validateMaxPlayers,
  validatePointsToPlayTo,
} from "./GameSettingsInput";
import { validate } from "./GameSettingsInputValidation";

vi.mock("./GameSettingsInput", () => {
  return {
    validateGamePassword: vi.fn(() => true),
    validateMaxGameRounds: vi.fn(() => true),
    validateMaxPlayers: vi.fn(() => true),
    validatePointsToPlayTo: vi.fn(() => true),
  };
});

describe("GameSettingsInput validate", () => {
  it("Should validate valid game settings", () => {
    const settings: Settings = {
      gamePassword: v4(),
      maxPlayers: 4,
      playingToPoints: 100,
      maxRounds: 10,
    };
    expect(validate(settings)).toBe(true);

    expect(validateGamePassword).toHaveBeenCalledWith(settings.gamePassword);
    expect(validateMaxGameRounds).toHaveBeenCalledWith(settings.maxRounds);
    expect(validateMaxPlayers).toHaveBeenCalledWith(settings.maxPlayers);
    expect(validatePointsToPlayTo).toHaveBeenCalledWith(
      settings.playingToPoints,
    );
  });

  it("Should not validate if game gamePassword is invalid", () => {
    (validateGamePassword as Mock).mockReturnValueOnce(false);
    const settings: Settings = {
      gamePassword: "invalid due to mock",
      maxPlayers: 4,
      playingToPoints: 100,
      maxRounds: 10,
    };
    expect(validate(settings)).toBe(false);
  });

  it("Should not validate if game maxPlayers is invalid", () => {
    (validateMaxPlayers as Mock).mockReturnValueOnce(false);
    const settings: Settings = {
      gamePassword: "valid",
      maxPlayers: 4,
      playingToPoints: 100,
      maxRounds: 10,
    };
    expect(validate(settings)).toBe(false);
  });

  it("Should not validate if game playingToPoints is invalid", () => {
    (validatePointsToPlayTo as Mock).mockReturnValueOnce(false);
    const settings: Settings = {
      gamePassword: "valid",
      maxPlayers: 4,
      playingToPoints: 100,
      maxRounds: 10,
    };
    expect(validate(settings)).toBe(false);
  });

  it("Should not validate if game maxRounds is invalid", () => {
    (validateMaxGameRounds as Mock).mockReturnValueOnce(false);
    const settings: Settings = {
      gamePassword: "valid",
      maxPlayers: 4,
      playingToPoints: 100,
      maxRounds: 10,
    };
    expect(validate(settings)).toBe(false);
  });
});
