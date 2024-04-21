import { GamePlayerList } from "../../gameState/gamePlayersList";
import RoundedWhite from "../containers/RoundedWhite";
import { For } from "solid-js";

interface Props {
  players: GamePlayerList;
}

export default function PlayerList(props: Readonly<Props>) {
  return (
    <RoundedWhite>
      <For each={props.players} fallback={<p>No players</p>}>
        {(player) => (
          <div class="flex flex-row gap-2">
            <p class="font-bold">{player.name}</p>
            <p class="font-mono">
              {player.connected ? "CONNECTED" : "DISCONNECTED"}
            </p>
            <p class="font-bold">{player.points} points</p>
          </div>
        )}
      </For>
    </RoundedWhite>
  );
}
