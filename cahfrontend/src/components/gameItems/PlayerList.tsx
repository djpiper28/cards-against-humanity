import { GamePlayerList } from "../../gameState/gamePlayersList";
import RoundedWhite from "../containers/RoundedWhite";
import { For, Show } from "solid-js";
import SubHeader from "../typography/SubHeader";
import { gameState } from "../../gameState/gameState";
import Button from "../buttons/Button";
import { PreviousWinner } from "../../gameLogicTypes";
import Card from "./Card";

interface Props {
  players: GamePlayerList;
  czarId: string;
  previousWinner: PreviousWinner;
}

export default function PlayerList(props: Readonly<Props>) {
  return (
    <RoundedWhite>
      <SubHeader text="Players" />
      <For each={props.players} fallback={<p>No players</p>}>
        {(player) => {
          const isWinner = () => player.id === props.previousWinner.playerId;

          return (
            <div
              class={`flex flex-col gap-3 ${isWinner() ? "border-2 rounded-2xl border-blue-800 p-2" : ""}`}
            >
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
                  <Show when={isWinner()}>
                    <p
                      class="font-bold text-blue-800"
                      id={`player-is-winner-${player.id}`}
                    >
                      WINNER
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
              <Show when={isWinner()}>
                <div class="flex flex-row gap-2">
                  <Card
                    id={props.previousWinner.blackCard?.id ?? 0}
                    isWhite={false}
                    cardText={props.previousWinner.blackCard?.bodyText ?? ""}
                    packName={`Play ${props.previousWinner.blackCard?.cardsToPlay}`}
                  />
                  <For each={props.previousWinner.whiteCards}>
                    {(card) => (
                      <Card
                        id={card?.id ?? 0}
                        cardText={card?.bodyText ?? ""}
                        packName=""
                        isWhite={true}
                      />
                    )}
                  </For>
                </div>
              </Show>
            </div>
          );
        }}
      </For>
    </RoundedWhite>
  );
}
