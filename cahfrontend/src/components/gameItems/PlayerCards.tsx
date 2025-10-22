import { For, Show } from "solid-js";
import Card from "./Card";
import Header from "../typography/Header";
import { gameState } from "../../gameState/gameState";

interface GameCard {
  id: string;
  name: string;
  pack: string;
}

interface Props {
  cards: GameCard[];
  selectedCardIds: string[];
  isJudging: boolean;
  isCzar: boolean;
  onSelectCard?: (cardId: string) => void;
}

export default function PlayerCards(props: Readonly<Props>) {
  const playerCardComp = () => (
    <div id="card-list" class="flex flex-row flex-wrap gap-2 justify-start">
      <For each={props.cards}>
        {(card) => {
          const isSelected = () =>
            !!props.selectedCardIds.find((x) => x === card.id);
          /**
           * Natural counting: 1 bound index.
           **/
          const cardNumber = () => props.selectedCardIds.indexOf(card.id) + 1;

          const id = `card-${card.id}`;
          const cardComp = () => (
            <>
              <Show when={isSelected()}>
                <div
                  class="absolute px-1 py-0.5 bg-blue-500 rounded-br-xl"
                  id={`selected-${cardNumber()}-${card.id}`}
                >
                  <p class="font-bold text-white">{cardNumber()}</p>
                </div>
              </Show>
              <Card
                id={card.id}
                isWhite={true}
                cardText={card.name}
                packName={card.pack}
              />
            </>
          );

          return (
            <>
              <Show when={props.isCzar}>{cardComp()}</Show>
              <Show when={props.czarId !== gameState.getPlayerId()}>
                <button
                  id={id}
                  class={
                    isSelected()
                      ? "border-4 border-blue-500 rounded-2xl bg-white cursor-pointer hover:scale-105"
                      : ""
                  }
                  onClick={() => props.onSelectCard?.(card.id)}
                >
                  {cardComp()}
                </button>
              </Show>
            </>
          );
        }}
      </For>
    </div>
  );

  const blockCardPlays = () => props.isCzar || props.isJudging;
  return (
    <>
      <Show when={blockCardPlays()}>
        <div class="relative flex w-fit">
          <div
            id="czar"
            class="absolute flex flex-col gap-3 top-0 left-0 right-0 bottom-0 z-10 justify-center items-center text-center bg-[#bababa60] rounded-2xl"
          >
            <Show when={props.isCzar}>
              <Header text="You are the Card Czar." />
            </Show>
            <Show when={props.isJudging}>
              <Header text="Judging in progress..." />
            </Show>
          </div>
          <div class="static blur">{playerCardComp()}</div>
        </div>
      </Show>
      <Show when={!blockCardPlays()}>{playerCardComp()}</Show>
    </>
  );
}
