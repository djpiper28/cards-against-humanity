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
          <div
            id="player-list"
            class="flex flex-row gap-3 ease-in duration-700 opacity-100"
          >
            <p class="font-bold" id={`player-name-${player.id}`}>
              {player.name}
            </p>
            <p class="font-mono" id={`player-status-${player.id}`}>
              {player.connected ? "CONNECTED" : "DISCONNECTED"}
            </p>
            <p class="font-bold" id={`player-points-${player.id}`}>
              {player.points} points
            </p>
          </div>
        )}
      </For>
    </RoundedWhite>
  );
}
