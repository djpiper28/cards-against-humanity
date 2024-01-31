import { useSearchParams } from "@solidjs/router";
import { gameIdParam, gameState, playerIdCookie } from "../gameState/gameState";
import { onMount } from "solid-js";
import { cookieStorage } from "@solid-primitives/storage";

export default function Join() {
  const [searchParams] = useSearchParams();

  onMount(() => {
    const gameId = searchParams[gameIdParam];
    const playerId = cookieStorage.getItem(playerIdCookie);

    gameState.setupState(gameId, playerId);
  });

  return <p>Joining the game...</p>;
}
