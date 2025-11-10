import { For, Show } from "solid-js";
import Card from "./Card";
import Header from "../typography/Header";

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

          return <Show when={props.isCzar}>{cardComp()}</Show>;
        }}
      </For>
    </div>
  );

  const blackCardJudging = () => props.isCzar || props.isJudging;
  return (
    <>
      <Show when={blackCardJudging()}>
        <div class="relative flex w-full">
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
          <Show when={!props.isCzar}>
            <div class="static blur">{playerCardComp()}</div>
          </Show>
        </div>
      </Show>
      <Show when={!blackCardJudging()}>{playerCardComp()}</Show>
    </>
  );
}
