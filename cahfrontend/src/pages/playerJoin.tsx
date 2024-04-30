import { createSignal } from "solid-js";
import Input, { InputType } from "../components/inputs/Input";
import { useNavigate, useSearchParams } from "@solidjs/router";
import { cookieStorage } from "@solid-primitives/storage";
import { validatePlayerName } from "../components/gameControls/GameSettingsInput";
import { apiClient, cookieOptions } from "../apiClient";
import { indexUrl, joinGameUrl } from "../routes";
import {
  authenticationCookie,
  gameIdParamCookie,
  gamePasswordCookie,
  gameState,
  playerIdCookie,
} from "../gameState/gameState";
import { MaxPasswordLength } from "../gameLogicTypes";
import RoundedWhite from "../components/containers/RoundedWhite";
import Header from "../components/typography/Header";
import Button from "../components/buttons/Button";
import clearGameCookies from "../gameState/clearGameCookies";

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

      <Button
        onClick={async () => {
          if (!validatePlayerName(playerName())) {
            return;
          }

          setError("");
          apiClient.games
            .joinCreate({
              gameId: gameId,
              playerName: playerName(),
              password: password(),
            })
            .then(async (res) => {
              try {
                await gameState.leaveGame();
                clearGameCookies();
              } catch (e) {
                console.log(e);
              }

              cookieStorage.setItem(
                gamePasswordCookie,
                password(),
                cookieOptions,
              );
              cookieStorage.setItem(
                playerIdCookie,
                res.data.playerId,
                cookieOptions,
              );
              cookieStorage.setItem(gameIdParamCookie, gameId, cookieOptions);
              cookieStorage.setItem(
                authenticationCookie,
                res.data.authentication,
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
      </Button>
    </RoundedWhite>
  );
}
