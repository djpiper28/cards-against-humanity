import { BlackCard, WhiteCard } from "../../gameLogicTypes";
import { For } from "solid-js";
import Card from "./Card";
import { gameState } from "../../gameState/gameState";

interface PlayerPlay {
  whiteCards: WhiteCard[];
  winner: boolean;
}

interface PlayerPlayProps {
  play: PlayerPlay;
  index: number;
  isCzar: boolean;
}

function PlayerCard(props: Readonly<PlayerPlayProps>) {
  const cardNode = (
    <div
      class={`flex flex-row gap-2 border-2 border-white rounded-2xl ${props.play.winner ? "border-blue-500" : ""}`}
    >
      {props.play.whiteCards.map((card, index) => (
        <Card
          id={props.index + index}
          cardText={card.bodyText}
          packName={`Player ${props.index + 1}`}
          isWhite={true}
        />
      ))}
    </div>
  );

  if (!props.isCzar) {
    return cardNode;
  }

  return (
    <button
      onClick={() => {
        gameState.czarSelectCards(props.play.whiteCards.map((x) => x.id));
      }}
    >
      {cardNode}
    </button>
  );
}

interface Props {
  blackCard: BlackCard;
  plays: PlayerPlay[];
  isCzar: boolean;
}

export default function CurrentRoundResults(props: Readonly<Props>) {
  return (
    <div class="flex flex-row gap-2 flex-wrap">
      <Card
        id="black-card"
        cardText={props.blackCard.bodyText}
        packName="Black card"
        isWhite={false}
      />
      <For each={props.plays}>
        {(play, index) => {
          return (
            <PlayerCard play={play} index={index()} isCzar={props.isCzar} />
          );
        }}
      </For>
    </div>
  );
}
