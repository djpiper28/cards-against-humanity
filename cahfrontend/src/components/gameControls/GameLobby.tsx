import { createSignal, onMount } from "solid-js";
import LoadingSlug from "../loading/LoadingSlug";
import { gameState } from "../../gameState/gameState";
import { GameStateInfo } from "../../gameLogicTypes";
import GameSettingsInput from "./GameSettingsInput";
import GameSettingsView from "./GameSettingsView";
import RoundedWhite from "../containers/RoundedWhite";

interface Props {
  state: GameStateInfo;
  gameOwner: boolean;
}

function GameLobbyLoaded(props: Readonly<Props>) {
  return (
    <RoundedWhite>
      <h1>{`${
        props.state.players.find((x) => x.id === props.state.gameOwnerId)?.name
      }'s Game`}</h1>
      {props.gameOwner ? (
        <GameSettingsInput
          settings={props.state.settings}
          setSettings={console.log}
          errorMessage="TODO: Implement me!"
        />
      ) : (
        <GameSettingsView settings={props.state.settings} />
      )}
    </RoundedWhite>
  );
}

export default function GameLobby() {
  const [state, setState] = createSignal<GameStateInfo | undefined>(undefined);
  const [gameOwner, setGameOwner] = createSignal(false);
  onMount(() => {
    gameState.onStateChange = (state?: GameStateInfo) => {
      console.log("State change detected");
      setState(state);

      if (gameState.isOwner()) {
        setGameOwner(true);
      }
    };
    // Just incase the update has already happened
    gameState.emitState();
  });

  return (
    <>
      {!!state() ? (
        <GameLobbyLoaded state={state()} gameOwner={gameOwner()} />
      ) : (
        <div class="flex flex-grow justify-center items-center text-2xl">
          Waiting for lobby information <LoadingSlug />
        </div>
      )}
    </>
  );
}
