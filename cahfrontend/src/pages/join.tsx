import { useNavigate, useSearchParams } from "@solidjs/router";
import { gameIdParam, gameState, playerIdCookie } from "../gameState/gameState";
import { onMount, createSignal } from "solid-js";
import { cookieStorage } from "@solid-primitives/storage";
import LoadingSlug from "../components/loading/LoadingSlug";
import GameLobby from "../components/gameControls/GameLobby";
import { indexUrl, playerJoinUrl } from "../routes";

export default function Join() {
  const [searchParams] = useSearchParams();
  const [connected, setConnected] = createSignal(false);
  const navigate = useNavigate();

  onMount(() => {
    const gameId = searchParams[gameIdParam];
    if (!gameId) {
      navigate("/");
      return;
    }

    const playerId = cookieStorage.getItem(playerIdCookie);
    if (!playerId) {
      navigate(`${playerJoinUrl}/${gameIdParam}=${gameId}`);
      return;
    }

    try {
      gameState.setupState(gameId, playerId);
      setConnected(true);
    } catch (e) {
      console.error(`Cannot setup the connection ${e}`);
      navigate(indexUrl);
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
