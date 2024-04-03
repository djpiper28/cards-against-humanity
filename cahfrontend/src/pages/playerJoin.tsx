import { createSignal } from "solid-js";
import Input, { InputType } from "../components/inputs/Input";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { cookieStorage } from "@solid-primitives/storage";
import { validatePlayerName } from "../components/gameControls/GameSettingsInput";
import { apiClient, cookieOptions } from "../apiClient";
import { indexUrl, joinGameUrl } from "../routes";
import {
  gameIdParamCookie,
  gamePasswordCookie,
  playerIdCookie,
} from "../gameState/gameState";
import { MaxPasswordLength } from "../gameLogicTypes";
import RoundedWhite from "../components/containers/RoundedWhite";
import Header from "../components/typography/Header";

export default function PlayerJoin() {
  const [searchParams] = useSearchParams();
  const [playerName, setPlayerName] = createSignal("");
  const [password, setPassword] = createSignal("");
  const [error, setError] = createSignal("");
  const naviagte = useNavigate();

  const gameId = searchParams[gameIdParamCookie];
  if (!gameId) {
    console.error("No gameId found, redirecting to index");
    naviagte(indexUrl);
    return;
  }

  return (
    <RoundedWhite>
      <Header text="Enter a username to join the game" />
      <Input
        inputType={InputType.Text}
        label="Player Name"
        value={playerName()}
        onChanged={(value) => setPlayerName(value)}
        placeHolder="John Smith"
        errorState={!validatePlayerName(playerName())}
      />
      <Input
        inputType={InputType.Text}
        autocomplete="password"
        label="Password, leave blank if none"
        value={password()}
        onChanged={(value) => setPassword(value)}
        placeHolder="Password"
        errorState={password().length > MaxPasswordLength}
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
              cookieStorage.setItem(playerIdCookie, res.data, cookieOptions);
              cookieStorage.setItem(gameIdParamCookie, gameId, cookieOptions);
              cookieStorage.setItem(
                gamePasswordCookie,
                password(),
                cookieOptions,
              );

              naviagte(
                `${joinGameUrl}?${gameIdParamCookie}=${encodeURIComponent(gameId)}`,
              );
            })
            .catch((res) => {
              setError(res.error.error);
            });
        }}
      >
        Join Game
      </button>
    </RoundedWhite>
  );
}
