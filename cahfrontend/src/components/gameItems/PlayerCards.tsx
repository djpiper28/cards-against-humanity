import { For } from "solid-js";
import Card from "./Card";

interface GameCard {
  id: string;
  name: string;
  pack: string;
}

interface Props {
  cards: GameCard[];
  selectedCardIds: string[];
  onSelectCard?: (cardId: string) => void;
}

export default function PlayerCards(props: Readonly<Props>) {
  return (
    <div class="flex flex-row flex-wrap gap-2 justify-between">
      <For each={props.cards}>
        {(card) => {
          const isSelected = () =>
            !!props.selectedCardIds.find((x) => x === card.id);
          /**
           * Natural counting: 1 bound index.
           **/
          const cardNumber = () => props.selectedCardIds.indexOf(card.id) + 1;
          return (
            <button
              class={
                isSelected()
                  ? "border-4 border-blue-500 rounded-2xl bg-white"
                  : ""
              }
              onClick={() => props.onSelectCard?.(card.id)}
            >
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
              <Card isWhite={true} cardText={card.name} packName={card.pack} />
            </button>
          );
        }}
      </For>
    </div>
  );
}
