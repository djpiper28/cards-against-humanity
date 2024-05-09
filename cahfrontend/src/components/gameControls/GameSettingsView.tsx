import { For } from "solid-js";
import { CardPack, GameSettings } from "../../gameLogicTypes";

interface Props {
  settings: GameSettings;
  packs: CardPack[];
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
              <span id={pack} class="font-bold">
                {props.packs.find((p) => p.id == pack)?.name}
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
