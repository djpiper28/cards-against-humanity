import { createSignal } from "solid-js";
import Input, { InputType } from "../components/inputs/Input";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { cookieStorage } from "@solid-primitives/storage";
import { validatePlayerName } from "../components/gameControls/GameSettingsInput";
import { apiClient } from "../apiClient";
import { indexUrl, joinGameUrl } from "../routes";
import { gameIdParam, playerIdCookie } from "../gameState/gameState";

export default function PlayerJoin() {
  const [searchParams] = useSearchParams();
  const [playerName, setPlayerName] = createSignal("");
  const [error, setError] = createSignal("");
  const naviagte = useNavigate();

  const gameId = searchParams[gameIdParam];
  if (!gameId) {
    console.error("No gameId found, redirecting to index");
    naviagte(indexUrl);
    return;
  }

  return (
    <>
      <h1 class="text-2xl font-bold">Enter a username to join the game</h1>
      <Input
        inputType={InputType.Text}
        label="Player Name"
        value={playerName()}
        onChanged={(value) => setPlayerName(value)}
        placeHolder="John Smith"
        errorState={!validatePlayerName(playerName())}
      />
      <p class="text-error-colour">{error()}</p>

      <button
        onClick={async () => {
          if (!validatePlayerName(playerName())) {
            return;
          }

          setError("");
          apiClient.games
            .joinCreate({
              gameId: gameId,
              playerName: playerName(),
            })
            .then((res) => {
              cookieStorage.setItem(playerIdCookie, res.data);
              cookieStorage.setItem(gameIdParam, gameId);
              naviagte(
                `${joinGameUrl}?${gameIdParam}=${encodeURIComponent(gameId)}`,
              );
            })
            .catch((res) => {
              setError(res.error.error);
            });
        }}
      >
        Join Game
      </button>
    </>
  );
}
