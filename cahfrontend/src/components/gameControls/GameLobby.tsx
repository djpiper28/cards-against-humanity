import { createSignal, onMount } from "solid-js";
import LoadingSlug from "../loading/LoadingSlug";
import { gameState } from "../../gameState/gameState";
import { GameStateInfo } from "../../gameLogicTypes";
import GameSettingsInput, { Settings } from "./GameSettingsInput";
import GameSettingsView from "./GameSettingsView";
import RoundedWhite from "../containers/RoundedWhite";
import Header from "../typography/Header";
import PlayerList from "../gameItems/PlayerList";
import { GamePlayerList } from "../../gameState/gamePlayersList";
import CardsSelector from "./CardsSelector";
import { apiClient } from "../../apiClient";
import { GameLogicCardPack } from "../../api";

interface Props {
  state: GameStateInfo;
  players: GamePlayerList;
  gameOwner: boolean;
  cardPacks: GameLogicCardPack[];
}

function GameLobbyLoaded(props: Readonly<Props>) {
  const [dirtyState, setDirtyState] = createSignal(false);
  const [selectedPackIds, setSelectedPackIds] = createSignal(
    props.state.settings.cardPacks.map((x) => x?.id ?? "no-id"),
  );
  const [gameSettings, setGameSettings] = createSignal<Settings>(
    props.state.settings,
  );

  return (
    <RoundedWhite>
      <Header
        text={`${
          props.state.players.find((x) => x.id === props.state.gameOwnerId)
            ?.name
        }'s Game`}
      />
      {props.gameOwner ? (
        <>
          <CardsSelector
            cards={props.cardPacks}
            selectedPackIds={selectedPackIds()}
            setSelectedPackIds={(packs) => {
              setDirtyState(true);
              setSelectedPackIds(packs);
            }}
          />
          <GameSettingsInput
            settings={gameSettings()}
            setSettings={(settings) => {
              setDirtyState(true);
              setGameSettings(settings);
            }}
            errorMessage="TODO: Implement me!"
          />
          <p id="settings-saved">
            {dirtyState() ? "Settings are NOT saved..." : "Settings are saved."}
          </p>
        </>
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

  const [packs, setPacks] = createSignal<GameLogicCardPack[]>([]);
  const [errorMessage, setErrorMessage] = createSignal("");

  onMount(async () => {
    try {
      const packs = await apiClient.res.packsList();
      const cardPacksList: GameLogicCardPack[] = [];
      const packData = packs.data;
      for (let cardId in packData) {
        cardPacksList.push({ ...packData[cardId] });
      }
      setPacks(
        cardPacksList.sort((a, b) => {
          if (!a.name || !b.name) return 0;
          return a.name.localeCompare(b.name);
        }),
      );
    } catch (err) {
      console.error(err);
      setErrorMessage(`Error getting card packs: ${err}`);
    }

    // Sync UI state with the source of truth
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
          cardPacks={packs()}
        />
      ) : (
        <div class="flex flex-grow justify-center items-center text-2xl">
          Waiting for lobby information <LoadingSlug />
        </div>
      )}
      <p id="error-message" class="text-error-colour">{errorMessage()}</p>
    </>
  );
}
