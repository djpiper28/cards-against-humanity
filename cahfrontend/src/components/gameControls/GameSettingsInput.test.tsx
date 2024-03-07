import { describe, expect, it } from "vitest";
import GameSettingsInput, { Settings } from "./GameSettingsInput";
import { screen, render, waitFor } from "solid-testing-library";

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
});
