import { createSignal, onMount } from "solid-js";
import LoadingSlug from "../loading/LoadingSlug";
import { gameState } from "../../gameState/gameState";
import { GameStateInfo } from "../../gameLogicTypes";
import GameSettingsInput from "./GameSettingsInput";
import GameSettingsView from "./GameSettingsView";
import RoundedWhite from "../containers/RoundedWhite";
import Header from "../typography/Header";
import PlayerList from "../gameItems/PlayerList";
import { GamePlayerList } from "../../gameState/gamePlayersList";

interface Props {
  state: GameStateInfo;
  players: GamePlayerList;
  gameOwner: boolean;
}

function GameLobbyLoaded(props: Readonly<Props>) {
  return (
    <RoundedWhite>
      <Header
        text={`${
          props.state.players.find((x) => x.id === props.state.gameOwnerId)
            ?.name
        }'s Game`}
      />
      {props.gameOwner ? (
        <GameSettingsInput
          settings={props.state.settings}
          setSettings={console.log}
          errorMessage="TODO: Implement me!"
        />
      ) : (
        <GameSettingsView settings={props.state.settings} />
      )}
      <PlayerList players={props.players} />
    </RoundedWhite>
  );
}

export default function GameLobby() {
  const [state, setState] = createSignal<GameStateInfo | undefined>(undefined);
  const [gameOwner, setGameOwner] = createSignal(false);
  const [players, setPlayers] = createSignal<GamePlayerList>([]);

  onMount(() => {
    gameState.onPlayerListChange = (players: GamePlayerList) => {
      setPlayers(players);
    };
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

  gameState.onPlayerListChange = (players: GamePlayerList) => {
    setPlayers(players);
  };
  gameState.onStateChange = (state?: GameStateInfo) => {
    console.log("State change detected");
    setState(state);

    if (gameState.isOwner()) {
      setGameOwner(true);
    }
  };

  return (
    <>
      {!!state() ? (
        <GameLobbyLoaded
          state={state()}
          gameOwner={gameOwner()}
          players={players()}
        />
      ) : (
        <div class="flex flex-grow justify-center items-center text-2xl">
          Waiting for lobby information <LoadingSlug />
        </div>
      )}
    </>
  );
}
