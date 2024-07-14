import { For } from "solid-js";
import Card from "./Card";
import Header from "../typography/Header";

interface GameCard {
  id: string;
  name: string;
  pack: string;
}

interface Props {
  isCzar: boolean;
  cards: GameCard[];
  selectedCardIds: string[];
  onSelectCard?: (cardId: string) => void;
}

export default function PlayerCards(props: Readonly<Props>) {
  const playerCardComp = (
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
              {isSelected() ? (
                <div
                  class="absolute px-1 py-0.5 bg-blue-500 rounded-br-xl"
                  id={`selected-${cardNumber()}-${card.id}`}
                >
                  <p class="font-bold text-white">{cardNumber()}</p>
                </div>
              ) : (
                <></>
              )}
              <Card
                id={card.id}
                isWhite={true}
                cardText={card.name}
                packName={card.pack}
              />
            </>
          );

          return props.isCzar ? (
            cardComp()
          ) : (
            <button
              id={id}
              class={
                isSelected()
                  ? "border-4 border-blue-500 rounded-2xl bg-white"
                  : ""
              }
              onClick={() => props.onSelectCard?.(card.id)}
            >
              {cardComp()}
            </button>
          );
        }}
      </For>
    </div>
  );

  if (props.isCzar) {
    return (
      <div class="relative flex w-fit">
        <div
          id="czar"
          class="absolute flex top-0 left-0 right-0 bottom-0 z-10 justify-center items-center text-center bg-[#aaaaaa30] rounded-2xl"
        >
          <Header text="You are the Card Czar" />
        </div>
        <div class="static">{playerCardComp}</div>
      </div>
    );
  } else {
    return playerCardComp;
  }
}
