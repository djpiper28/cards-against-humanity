import { useNavigate, useSearchParams } from "@solidjs/router";
import { gameIdParam, gameState, playerIdCookie } from "../gameState/gameState";
import { onMount, createSignal } from "solid-js";
import { cookieStorage } from "@solid-primitives/storage";
import LoadingSlug from "../components/loading/LoadingSlug";

export default function Join() {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  onMount(() => {
    const gameId = searchParams[gameIdParam];
    if (!gameId) {
      navigate("/");
      return;
    }

    const playerId = cookieStorage.getItem(playerIdCookie);
    if (!playerId) {
      navigate(`/player-join?${gameIdParam}=${gameId}`);
      return;
    }

    gameState.setupState(gameId, playerId);
  });

  return (
    <div class="flex flex-grow justify-center items-center text-2xl">
      Joining the game <LoadingSlug />
    </div>
  );
}
