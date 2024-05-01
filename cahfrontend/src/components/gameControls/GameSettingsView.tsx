import { For } from "solid-js";
import { GameSettings } from "../../gameLogicTypes";

interface Props {
  settings: GameSettings;
}

export default function GameSettingsView(props: Readonly<Props>) {
  return (
    <div>
      <h2>Game Settings</h2>
      <div>
        <p>
          Game Password:{" "}
          <span id="game-password">{props.settings.gamePassword}</span>
        </p>
        <p>
          Max Players: <span id="max-players">{props.settings.maxPlayers}</span>
        </p>
        <p>
          Points To Play To:{" "}
          <span id="playing-to-points">{props.settings.playingToPoints}</span>
        </p>
        <p>
          Max Rounds:{" "}
          <span id="max-game-rounds">{props.settings.maxRounds}</span>
        </p>
      </div>
      <h3>Card Packs</h3>
      <span>
        <For each={props.settings.cardPacks}>
          {(pack, i) => (
            <>
              <span id={pack?.id} class="font-bold">
                {pack?.name}
              </span>
              {i() == props.settings.cardPacks.length - 1 ? (
                ""
              ) : (
                <span class="pr-2">,</span>
              )}
            </>
          )}
        </For>
      </span>
    </div>
  );
}
