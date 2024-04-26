import { Show, createEffect, createSignal, onMount } from "solid-js";
import LoadingSlug from "../loading/LoadingSlug";
import { gameIdParamCookie, gamePasswordCookie, gameState, playerIdCookie } from "../../gameState/gameState";
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
import Button from "../buttons/Button";
import { cookieStorage } from "@solid-primitives/storage";
import { useNavigate } from "@solidjs/router";
import { indexUrl } from "../../routes";

interface LobbyLoadedProps {
  setSettings: (settings: Settings) => void;
  setSelectedPackIds: (ids: string[]) => void;
  players: GamePlayerList;
  commandError: string;
  setCommandError: (error: string) => void;
  dirtyState: boolean;
  cardPacks: GameLogicCardPack[];
  state: GameStateInfo;
  setStateAsDirty: () => void;
}

function GameLobbyLoaded(props: Readonly<LobbyLoadedProps>) {
  const updateSettings = () => {
    const changeSettings: RpcChangeSettingsMsg = {
      settings: {
        ...props.state.settings,
        cardPacks: props.state.settings.cardPacks,
      },
    };
    gameState.sendRpcMessage(MsgChangeSettings, changeSettings);
  };

  const state = () => props.state;
  const isGameOwner = () => state().gameOwnerId === gameState.getPlayerId();
  const settings = () => state().settings;
  const dirtyState = () => props.dirtyState;
  const navigate = useNavigate();

  createEffect(() => {
    cookieStorage.setItem(gamePasswordCookie, settings().gamePassword);
  });

  return (
    <RoundedWhite>
      <div class="flex flex-row justify-between flex-wrap">
        <Header
          text={`${
            props.players.find((x) => x.id === props.state.gameOwnerId)?.name
          }'s Game`}
        />
        <Button
          id="leave-game"
          onClick={() => {
            gameState
              .leaveGame()
              .then(() => console.log("Left game successfully"))
              .catch(console.error);
            navigate(indexUrl);
            cookieStorage.removeItem(gamePasswordCookie);
            cookieStorage.removeItem(gameIdParamCookie);
            cookieStorage.removeItem(playerIdCookie);
          }}
        >
          Leave Game
        </Button>
      </div>
      <Show when={isGameOwner()}>
        <Header text="Change your game's settings" />
        <CardsSelector
          cards={props.cardPacks}
          selectedPackIds={settings().cardPacks.map((x) => x?.id ?? "no-id")}
          setSelectedPackIds={(packs) => {
            props.setStateAsDirty();
            props.setSelectedPackIds(packs);
            updateSettings();
          }}
        />
        <GameSettingsInput
          settings={settings()}
          setSettings={(settings) => {
            props.setStateAsDirty();
            props.setSettings(settings);
            updateSettings();
          }}
          errorMessage={props.commandError}
        />
        <p id="settings-saved">
          <Show when={dirtyState()} fallback={"Settings are saved."}>
            <div class="flex flex-row gap-2">
              <p class="text-error-colour">Settings are NOT saved.</p>
              <Button onClick={updateSettings}>Save</Button>
              <Button
                onClick={() => {
                  gameState.emitState();
                  props.setCommandError("");
                }}
              >
                Reset
              </Button>
            </div>
          </Show>
        </p>
      </Show>

      <Show when={!isGameOwner()}>
        <GameSettingsView settings={settings()} />
      </Show>

      <Show when={isGameOwner()}>
        <Show
          when={!dirtyState()}
          fallback={"Cannot start a game with unsaved changes."}
        >
          <Button onClick={() => {}}>Start Game</Button>
        </Show>
      </Show>
      <PlayerList players={props.players} />
    </RoundedWhite>
  );
}

const emptyState: GameStateInfo = {
  id: "",
  gameOwnerId: "",
  settings: {
    cardPacks: [],
    maxPlayers: 0,
    maxRounds: 0,
    playingToPoints: 0,
    gamePassword: "",
  },
  currentRound: 0,
  creationTime: "",
  gameState: 0,
  players: [],
  currentCardCzarId: "",
};

export default function GameLobby() {
  const [state, setState] = createSignal<GameStateInfo>(emptyState);
  const [players, setPlayers] = createSignal<GamePlayerList>([]);

  const [packs, setPacks] = createSignal<GameLogicCardPack[]>([]);
  const [errorMessage, setErrorMessage] = createSignal("");
  const [commandError, setCommandError] = createSignal("");

  const [dirtyState, setDirtyState] = createSignal(false);

  const setupHandlers = () => {
    // Sync UI state with the source of truth
    gameState.onPlayerListChange = (players: GamePlayerList) => {
      setPlayers(players);
    };
    gameState.onStateChange = (state?: GameStateInfo) => {
      console.log("State change detected");
      setDirtyState(false);
      setState(state ?? emptyState);
    };
    gameState.onChangeSettings = (settings: RpcChangeSettingsMsg) => {
      setDirtyState(false);

      const newState = state() ?? emptyState;
      const data = settings.settings as GameSettings;

      newState.settings.maxRounds = data.maxRounds;
      newState.settings.maxPlayers = data.maxPlayers;
      newState.settings.playingToPoints = data.playingToPoints;
      newState.settings.gamePassword = data.gamePassword;

      setState(newState);
    };
    gameState.onCommandError = (error: RpcCommandErrorMsg) => {
      setCommandError(error.reason);
    };
  };

  setupHandlers();
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
        })
      );
    } catch (err) {
      console.error(err);
      setErrorMessage(`Error getting card packs: ${err}`);
    }

    setupHandlers();
    // Just incase the update has already happened
    gameState.emitState();
  });

  return (
    <>
      <Show when={!!state()}>
        <GameLobbyLoaded
          state={state()}
          setSettings={(settings) => {
            const newState = state();
            newState.settings.maxRounds = settings.maxRounds;
            newState.settings.maxPlayers = settings.maxPlayers;
            newState.settings.playingToPoints = settings.playingToPoints;
            newState.settings.gamePassword = settings.gamePassword;
            setState(newState);
          }}
          setSelectedPackIds={(ids) => {
            const newState = state();
            newState.settings.cardPacks = ids.map((id) =>
              packs().find((x) => x.id === id)
            );
            setState(newState);
          }}
          players={players()}
          cardPacks={packs()}
          commandError={commandError()}
          setCommandError={setCommandError}
          dirtyState={dirtyState()}
          setStateAsDirty={() => setDirtyState(true)}
        />
      </Show>

      <Show when={!state()}>
        <div class="flex flex-grow justify-center items-center text-2xl">
          Waiting for lobby information <LoadingSlug />
        </div>
      </Show>

      <p id="error-message" class="text-error-colour">
        {errorMessage()}
      </p>
    </>
  );
}
