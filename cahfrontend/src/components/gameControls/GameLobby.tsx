import { createSignal, onMount, runWithOwner } from "solid-js";
import LoadingSlug from "../loading/LoadingSlug";
import { gameState } from "../../gameState/gameState";
import { GameSettings, GameStateInfo } from "../../gameLogicTypes";
import GameSettingsInput, { Settings } from "./GameSettingsInput";
import GameSettingsView from "./GameSettingsView";
import RoundedWhite from "../containers/RoundedWhite";
import Header from "../typography/Header";
import PlayerList from "../gameItems/PlayerList";
import { GamePlayerList } from "../../gameState/gamePlayersList";
import CardsSelector from "./CardsSelector";
import { apiClient } from "../../apiClient";
import { GameLogicCardPack } from "../../api";
import {
  MsgChangeSettings,
  RpcChangeSettingsMsg,
  RpcCommandErrorMsg,
} from "../../rpcTypes";

interface Props {
  state: GameStateInfo;
  players: GamePlayerList;
  gameOwner: boolean;
  cardPacks: GameLogicCardPack[];
  commandError: string;
  dirtyState: boolean;
  setStateAsDirty: () => void;
}

function GameLobbyLoaded(props: Readonly<Props>) {
  const [selectedPackIds, setSelectedPackIds] = createSignal(
    props.state.settings.cardPacks.map((x) => x?.id ?? "no-id"),
  );
  const [gameSettings, setGameSettings] = createSignal<Settings>(
    props.state.settings,
  );

  const updateSettings = () => {
    const changeSettings: RpcChangeSettingsMsg = {
      settings: {
        ...gameSettings(),
        cardPacks: selectedPackIds().map((id) =>
          props.cardPacks.find((x) => x.id === id),
        ),
      },
    };

    gameState.sendRpcMessage(MsgChangeSettings, changeSettings);
  };

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
          <Header text="Change your game's settings" />
          <CardsSelector
            cards={props.cardPacks}
            selectedPackIds={selectedPackIds()}
            setSelectedPackIds={(packs) => {
              props.setStateAsDirty();
              setSelectedPackIds(packs);
              updateSettings();
            }}
          />
          <GameSettingsInput
            settings={gameSettings()}
            setSettings={(settings) => {
              props.setStateAsDirty();
              setGameSettings(settings);
              updateSettings();
            }}
            errorMessage={props.commandError}
          />
          <p id="settings-saved">
            {props.dirtyState ? (
              <>
                Settings are NOT saved <LoadingSlug />{" "}
              </>
            ) : (
              "Settings are saved."
            )}
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
  const [commandError, setCommandError] = createSignal("");

  const [dirtyState, setDirtyState] = createSignal(false);

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
      setDirtyState(false);
      setState(state);

      if (gameState.isOwner()) {
        setGameOwner(true);
      }
    };
    gameState.onChangeSettings = (settings: RpcChangeSettingsMsg) => {
      setDirtyState(false);

      const newState = state();
      newState.settings = settings.settings;
      setState(newState);
    };
    gameState.onCommandError = (error: RpcCommandErrorMsg) => {
      setCommandError(error.reason);
    };

    // Just incase the update has already happened
    gameState.emitState();
  });

  return (
    <>
      {!!state() ? (
        <GameLobbyLoaded
          state={state()}
          gameOwner={gameOwner()}
          players={players()}
          cardPacks={packs()}
          commandError={commandError()}
          dirtyState={dirtyState()}
          setStateAsDirty={() => setDirtyState(true)}
        />
      ) : (
        <div class="flex flex-grow justify-center items-center text-2xl">
          Waiting for lobby information <LoadingSlug />
        </div>
      )}
      <p id="error-message" class="text-error-colour">
        {errorMessage()}
      </p>
    </>
  );
}
