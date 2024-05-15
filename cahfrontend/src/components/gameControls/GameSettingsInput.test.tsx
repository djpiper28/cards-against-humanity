import { describe, expect, it } from "vitest";
import GameSettingsInput, {
  Settings,
  validateGamePassword,
  validateMaxGameRounds,
  validateMaxPlayers,
  validatePointsToPlayTo,
  validatePlayerName,
} from "./GameSettingsInput";
import { screen, render, waitFor } from "@solidjs/testing-library";
import {
  MaxPasswordLength,
  MaxPlayers,
  MaxPlayingToPoints,
  MaxRounds,
  MinPlayerNameLength,
  MinPlayers,
  MinPlayingToPoints,
  MinRounds,
  MaxPlayerNameLength,
} from "../../gameLogicTypes";

describe("GameSettingsInput", () => {
  it("Should render with the settings", async () => {
    const settings: Settings = {
      gamePassword: "password",
      maxPlayers: 4,
      playingToPoints: 100,
      maxRounds: 10,
    };

    render(() => (
      <GameSettingsInput
        settings={settings}
        setSettings={() => {}}
        errorMessage=""
      />
    ));

    waitFor(async () => {
      expect((await screen.findByLabelText("Game Password")).nodeValue).toBe(
        settings.gamePassword,
      );
      expect((await screen.findByLabelText("Max Players")).nodeValue).toBe(
        settings.maxPlayers.toString(),
      );
      expect(
        (await screen.findByLabelText("Points To Play To")).nodeValue,
      ).toBe(settings.playingToPoints.toString());
      expect((await screen.findByLabelText("Max Game Rounds")).nodeValue).toBe(
        settings.maxPlayers.toString(),
      );
    });
  });

  it("Should validate the game password", () => {
    expect(validateGamePassword("")).toBe(true);
    expect(validateGamePassword("password")).toBe(true);
    expect(validateGamePassword("*".repeat(MaxPasswordLength))).toBe(true);
    expect(validateGamePassword("*".repeat(MaxPasswordLength + 1))).toBe(false);
  });

  it("Should validate the max players", () => {
    expect(validateMaxPlayers(-1)).toBe(false);
    expect(validateMaxPlayers(0)).toBe(false);
    expect(validateMaxPlayers(MinPlayers - 1)).toBe(false);
    expect(validateMaxPlayers(MinPlayers)).toBe(true);
    expect(validateMaxPlayers(MinPlayers + 1)).toBe(true);
    expect(validateMaxPlayers(MaxPlayers)).toBe(true);
    expect(validateMaxPlayers(MaxPlayers + 1)).toBe(false);
  });

  it("Should validate the points to play to", () => {
    expect(validatePointsToPlayTo(-1)).toBe(false);
    expect(validatePointsToPlayTo(0)).toBe(false);
    expect(validatePointsToPlayTo(MinPlayingToPoints - 1)).toBe(false);
    expect(validatePointsToPlayTo(MinPlayingToPoints)).toBe(true);
    expect(validatePointsToPlayTo(MinPlayingToPoints + 1)).toBe(true);
    expect(validatePointsToPlayTo(MaxPlayingToPoints)).toBe(true);
    expect(validatePointsToPlayTo(MaxPlayingToPoints + 1)).toBe(false);
  });

  it("Should validate the max game rounds", () => {
    expect(validateMaxGameRounds(-1)).toBe(false);
    expect(validateMaxGameRounds(0)).toBe(false);
    expect(validateMaxGameRounds(MinRounds - 1)).toBe(false);
    expect(validateMaxGameRounds(MinRounds)).toBe(true);
    expect(validateMaxGameRounds(MinRounds + 1)).toBe(true);
    expect(validateMaxGameRounds(MaxRounds)).toBe(true);
    expect(validateMaxGameRounds(MaxRounds + 1)).toBe(false);
  });

  it("Should validate the player name", () => {
    expect(validatePlayerName("")).toBe(false);
    expect(validatePlayerName("a")).toBe(false);
    expect(validatePlayerName("a".repeat(MinPlayerNameLength - 1))).toBe(false);
    expect(validatePlayerName("a".repeat(MinPlayerNameLength))).toBe(true);
    expect(validatePlayerName("a".repeat(MinPlayerNameLength + 1))).toBe(true);
    expect(validatePlayerName("a".repeat(MaxPlayerNameLength))).toBe(true);
    expect(validatePlayerName("a".repeat(MaxPlayerNameLength + 1))).toBe(false);
  });
});
