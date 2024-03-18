import { createSignal, onMount } from "solid-js";
import LoadingSlug from "../loading/LoadingSlug";
import { gameState } from "../../gameState/gameState";
import { GameStateInfo } from "../../gameLogicTypes";

export default function GameLobby() {
  const [state, setState] = createSignal<GameStateInfo | undefined>(undefined);
  gameState.onStateChange = (state?: GameStateInfo) => {
    setState(state);
  };

  onMount(() => {
    // Just incase the update has already happened
    gameState.emitState();
  });

  return (
    <>
      {state() ? (
        <h1>LOADED</h1>
      ) : (
        <div class="flex flex-grow justify-center items-center text-2xl">
          Waiting for lobby information <LoadingSlug />
        </div>
      )}
    </>
  );
}
