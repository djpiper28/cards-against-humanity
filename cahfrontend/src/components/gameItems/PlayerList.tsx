import { GamePlayerList } from "../../gameState/gamePlayersList";
import RoundedWhite from "../containers/RoundedWhite";
import { For } from "solid-js";
import SubHeader from "../typography/SubHeader";

interface Props {
  players: GamePlayerList;
}

export default function PlayerList(props: Readonly<Props>) {
  return (
    <RoundedWhite>
      <SubHeader text="Players" />
      <For each={props.players} fallback={<p>No players</p>}>
        {(player) => (
          <div class="flex flex-row gap-3">
            <p class="font-bold" id={`${player.id}-player-name`}>
              {player.name}
            </p>
            <p class="font-mono" id={`${player.id}-player-status`}>
              {player.connected ? "CONNECTED" : "DISCONNECTED"}
            </p>
            <p class="font-bold" id={`${player.id}-player-points`}>
              {player.points} points
            </p>
          </div>
        )}
      </For>
    </RoundedWhite>
  );
}
