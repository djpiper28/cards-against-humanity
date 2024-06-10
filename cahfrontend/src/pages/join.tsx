import { useNavigate, useSearchParams } from "@solidjs/router";
import {
  gameIdParamCookie,
  gamePasswordCookie,
  gameState,
  playerIdCookie,
} from "../gameState/gameState";
import { onMount, createSignal } from "solid-js";
import { cookieStorage } from "@solid-primitives/storage";
import LoadingSlug from "../components/loading/LoadingSlug";
import GameLobby from "../components/gameControls/GameLobby";
import {
  gameErrorUrl as joinErrorUrl,
  indexUrl,
  playerJoinUrl,
  gameErrorUrl,
  joinGameUrl,
} from "../routes";
import { errorMessageParam } from "./gameError";

export default function Join() {
  const [searchParams] = useSearchParams();
  const [connected, setConnected] = createSignal(false);
  const navigate = useNavigate();

  onMount(() => {
    const gameId = searchParams[gameIdParamCookie];
    const alternateGameId = cookieStorage.getItem(gameIdParamCookie);
    if (!gameId && alternateGameId) {
      navigate(`${joinGameUrl}?${gameIdParamCookie}=${alternateGameId}`);
      return;
    }

    if (!gameId) {
      console.error("There is no gameId, redirecting to index");
      navigate(indexUrl);
      return;
    }

    const playerId = cookieStorage.getItem(playerIdCookie);
    const cookieGameId = cookieStorage.getItem(gameIdParamCookie);
    const password = cookieStorage.getItem(gamePasswordCookie) ?? "";
    if (!playerId) {
      console.log("No playerId found, redirecting to player join page");
      navigate(`${playerJoinUrl}?${gameIdParamCookie}=${gameId}`);
      return;
    }

    if (!cookieGameId) {
      console.log("There is no gameId cookie");
      navigate(`${playerJoinUrl}?${gameIdParamCookie}=${gameId}`);
      return;
    }

    if (cookieGameId && cookieGameId !== gameId) {
      console.log(
        "Stored playerId is for a different game, redirecting to player join page",
      );
      navigate(`${playerJoinUrl}?${gameIdParamCookie}=${gameId}`);
      return;
    }

    try {
      gameState.setupState(gameId, playerId, password);
      gameState.onError = (msg: string) => {
        navigate(
          `${gameErrorUrl}?${errorMessageParam}=${encodeURIComponent(msg)}`,
        );
      };
      setConnected(true);
    } catch (e) {
      console.error(`Cannot setup the connection ${e}`);
      navigate(`${joinErrorUrl}?${gameIdParamCookie}=${gameId}`);
    }
  });

  return (
    <>
      {connected() ? (
        <GameLobby />
      ) : (
        <div class="flex flex-grow justify-center items-center text-2xl">
          Connecting to the game server <LoadingSlug />
        </div>
      )}
    </>
  );
}
