import { For } from "solid-js";
import Card from "./Card";

interface GameCard {
  id: string;
  name: string;
  pack: string;
  played: boolean;
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
        {(card) => (
          <button
            class={
              !!props.selectedCardIds.find((x) => x === card.id)
                ? "border-4 border-blue-500 rounded-2xl bg-white"
                : ""
            }
            onClick={() => props.onSelectCard?.(card.id)}
          >
            <Card
              isWhite={!card.played}
              cardText={card.name}
              packName={card.pack}
            />
          </button>
        )}
      </For>
    </div>
  );
}
