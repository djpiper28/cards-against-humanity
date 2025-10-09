import { GamePlayerList } from "../../gameState/gamePlayersList";
import RoundedWhite from "../containers/RoundedWhite";
import { For, Show } from "solid-js";
import SubHeader from "../typography/SubHeader";
import { gameState } from "../../gameState/gameState";
import Button from "../buttons/Button";

interface Props {
  players: GamePlayerList;
  czarId: string;
}

export default function PlayerList(props: Readonly<Props>) {
  return (
    <RoundedWhite>
      <SubHeader text="Players" />
      <For each={props.players} fallback={<p>No players</p>}>
        {(player) => (
          <div class="flex flex-row flex-wrap gap-3 justify-between">
            <div
              id="player-list"
              class="flex flex-row gap-3 ease-in duration-700 opacity-100"
            >
              <p class="font-bold" id={`player-name-${player.id}`}>
                {player.name}
              </p>
              <p
                class={`font-mono ${player.connected ? "text-black" : "text-red-700 font-bold"}`}
                id={`player-status-${player.id}`}
              >
                {player.connected ? "CONNECTED" : "DISCONNECTED"}
              </p>
              <p class="font-bold" id={`player-points-${player.id}`}>
                {player.points} points
              </p>
              <Show when={player.hasPlayed}>
                <p class="font-bold" id={`player-has-played-${player.id}`}>
                  HAS PLAYED
                </p>
              </Show>
              <Show when={player.id === props.czarId}>
                <p class="font-bold" id={`player-has-played-${player.id}`}>
                  CARD CZAR
                </p>
              </Show>
              <Show when={player.id === gameState.getPlayerId()}>
                <p class="font-light" id="player-you">
                  (you)
                </p>
              </Show>
            </div>
            <Show
              when={
                gameState.isOwner() && player.id !== gameState.getPlayerId()
              }
            >
              <Button
                id={`kick-player-${player.id}`}
                onClick={() => {
                  gameState.czarKickPlayer(player.id);
                }}
              >
                Kick {player.name}
              </Button>
            </Show>
          </div>
        )}
      </For>
    </RoundedWhite>
  );
}
