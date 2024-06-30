import { LobbyLoadedProps } from "./gameLoadedProps";
import { Show } from "solid-js";
import { gameState } from "../../gameState/gameState";
import { MsgChangeSettings, RpcChangeSettingsMsg } from "../../rpcTypes";
import CardsSelector from "./CardsSelector";
import GameSettingsInput from "./GameSettingsInput";
import GameSettingsView from "./GameSettingsView";
import Button from "../buttons/Button";

export function GameNotStartedView(props: Readonly<LobbyLoadedProps>) {
  const updateSettings = () => {
    const changeSettings: RpcChangeSettingsMsg = {
      settings: {
        ...props.state.settings,
        cardPacks: props.state.settings.cardPacks,
      },
    };
    gameState.sendRpcMessage(MsgChangeSettings, changeSettings);
  };

  const isGameOwner = () => props.state.ownerId === gameState.getPlayerId();
  const settings = () => props.state.settings;
  const dirtyState = () => props.dirtyState;

  return (
    <>
      <Show when={isGameOwner()}>
        <CardsSelector
          cards={props.cardPacks}
          selectedPackIds={settings().cardPacks}
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
              <Button id="save-settings" onClick={updateSettings}>
                Save
              </Button>
              <Button
                id="reset-settings"
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
        <GameSettingsView settings={settings()} packs={props.cardPacks} />
      </Show>

      <Show when={isGameOwner()}>
        <Show
          when={!dirtyState()}
          fallback={"Cannot start a game with unsaved changes."}
        >
          <Button
            id="start-game"
            onClick={async () => {
              try {
                await gameState.startGame();
              } catch (e) {
                props.setCommandError(
                  "Unable to start game. Please try again.",
                );
              }
            }}
          >
            Start Game
          </Button>
        </Show>
      </Show>
    </>
  );
}
