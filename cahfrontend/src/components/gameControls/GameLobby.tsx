import { Show, createEffect, createSignal, onMount } from "solid-js";
import LoadingSlug from "../loading/LoadingSlug";
import { gamePasswordCookie, gameState } from "../../gameState/gameState";
import { GameSettings, GameStateInLobby } from "../../gameLogicTypes";
import RoundedWhite from "../containers/RoundedWhite";
import Header from "../typography/Header";
import PlayerList from "../gameItems/PlayerList";
import { GamePlayerList } from "../../gameState/gamePlayersList";
import { apiClient } from "../../apiClient";
import { GameLogicCardPack } from "../../api";
import {
  RpcChangeSettingsMsg,
  RpcCommandErrorMsg,
  RpcRoundInformationMsg,
} from "../../rpcTypes";
import Button from "../buttons/Button";
import { cookieStorage } from "@solid-primitives/storage";
import { useNavigate } from "@solidjs/router";
import { indexUrl } from "../../routes";
import clearGameCookies from "../../gameState/clearGameCookies";
import SubHeader from "../typography/SubHeader";
import { GameLobbyState } from "../../gameState/gameLobbyState";
import PlayerCards from "../gameItems/PlayerCards";
import { LobbyLoadedProps } from "./gameLoadedProps";
import { GameNotStartedView } from "./GameNotStartedView";
import Card from "../gameItems/Card";

// Exported for testing reasons
export function GameLobbyLoaded(props: Readonly<LobbyLoadedProps>) {
  const isGameOwner = () => props.state.ownerId === gameState.getPlayerId();
  const settings = () => props.state.settings;
  const isGameStarted = () => props.state.gameState !== GameStateInLobby;

  createEffect(() => {
    cookieStorage.setItem(gamePasswordCookie, settings().gamePassword);
  });

  return (
    <RoundedWhite>
      <div class="flex flex-row gap-3 justify-between flex-wrap">
        <div class="flex flex-row flex-wrap gap-5 w-fit items-center">
          <Header
            text={`${
              props.players.find((x) => x.id === props.state.ownerId)?.name
            }'s Game`}
          />
          <Show when={isGameOwner() && !isGameStarted()}>
            <SubHeader text="Change your game's settings" />
          </Show>
        </div>
        <Button
          id="leave-game"
          onClick={() => {
            gameState
              .leaveGame()
              .then(() => {
                clearGameCookies();
                console.log("Left game successfully");
                props.navigate(indexUrl);
              })
              .catch((e) => {
                console.error(e);
                props.setCommandError(
                  "Unable to leave game. Please try again.",
                );
              });
          }}
        >
          Leave Game
        </Button>
      </div>

      <div class="flex flex-col gap-3">
        <Show when={!isGameStarted()}>
          <GameNotStartedView {...props} />
        </Show>

        <Show when={isGameStarted() && props.roundState}>
          <Card
            isWhite={false}
            cardText={props.roundState.blackCard.bodyText}
            packName={`${props.roundState.blackCard.cardsToPlay} white cards to play`}
          />
          <PlayerCards
            cards={props.roundState.yourHand.map((x) => {
              return {
                id: x.id.toString(),
                name: x.bodyText,
                pack: "Your hand",
              };
            })}
            selectedCardIds={props.roundState.yourPlays.map((x) =>
              x.id.toString(),
            )}
            onSelectCard={(id) =>
              props.setSelectedCardIds([
                ...props.roundState.yourPlays.map((x) => x.id.toString()),
                id,
              ])
            }
          />
        </Show>

        <PlayerList players={props.players} />
      </div>
    </RoundedWhite>
  );
}

const emptyState: GameLobbyState = {
  ownerId: "",
  settings: {
    gamePassword: "",
    maxPlayers: 0,
    cardPacks: [],
    maxRounds: 0,
    playingToPoints: 0,
  },
  creationTime: new Date(),
  gameState: 0,
};

export default function GameLobby() {
  const [state, setState] = createSignal<GameLobbyState | undefined>(undefined);
  const [roundState, setRoundState] = createSignal<
    RpcRoundInformationMsg | undefined
  >(undefined);
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
    gameState.onLobbyStateChange = (state?: GameLobbyState) => {
      console.log("State change detected");
      setDirtyState(false);
      setState(state ?? emptyState);
    };
    gameState.onChangeSettings = (settings: RpcChangeSettingsMsg) => {
      setDirtyState(false);

      const newState = structuredClone(state()!);
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
    gameState.onRoundStateChange = (state?: RpcRoundInformationMsg) => {
      setRoundState(state);
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
        }),
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
          state={state()!}
          roundState={roundState()!}
          setSettings={(settings) => {
            const newState = state()!;
            newState.settings.maxRounds = settings.maxRounds;
            newState.settings.maxPlayers = settings.maxPlayers;
            newState.settings.playingToPoints = settings.playingToPoints;
            newState.settings.gamePassword = settings.gamePassword;
            setState(newState);
          }}
          setSelectedPackIds={(ids) => {
            const newState = state()!;
            newState.settings.cardPacks = ids;
            setState(newState);
          }}
          setSelectedCardIds={async (ids) => {
            const rState = roundState()!;
            if (ids.length > rState.blackCard.cardsToPlay) {
              ids = ids.slice(1);
            }

            const newRoundState = structuredClone(rState);
            newRoundState.yourPlays = ids.map((id) => {
              return rState.yourHand.find((x) => x.id.toString() === id)!;
            });
            setRoundState(newRoundState);

            // Do last not to block the UI updates
            if (ids.length === rState.blackCard.cardsToPlay) {
              gameState.playCards(ids.map((x) => parseInt(x)));
            }
          }}
          players={players()}
          cardPacks={packs()}
          commandError={commandError()}
          setCommandError={setCommandError}
          dirtyState={dirtyState()}
          setStateAsDirty={() => setDirtyState(true)}
          navigate={useNavigate()}
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
