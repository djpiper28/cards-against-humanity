import Card from "./Card";

interface WhiteCards {
  id: string;
  name: string;
  pack: string;
}

interface PlayerPlay {
  whiteCards: WhiteCards[];
  winner: boolean;
}

interface BlackCard {
  name: string;
  pack: string;
  cardsToPlay: number;
}

interface Props {
  blackCard: BlackCard;
  plays: PlayerPlay[];
  showCards: boolean;
}

interface PlayerPlayProps {
  play: PlayerPlay;
  showCards: boolean;
  index: number;
}

function PlayerCard(props: Readonly<PlayerPlayProps>) {
  return (
    <div
      class={`flex flex-row gap-2 border-2 border-white rounded-2xl ${props.play.winner ? "border-blue-500" : ""}`}
    >
      {props.play.whiteCards.map((card, index) => (
        <Card
          id={props.index + index}
          cardText={props.showCards ? card.name : "Cards Against Humanity"}
          packName={props.showCards ? card.pack : ""}
          isWhite={true}
        />
      ))}
    </div>
  );
}

export default function CurrentRoundResults(props: Readonly<Props>) {
  return (
    <div class="flex flex-row gap-2">
      <Card
        id="black-card"
        cardText={props.blackCard.name}
        packName={props.blackCard.pack}
        isWhite={false}
      />
      <div class="flex flex-row flex-wrap gap-2">
        {props.plays.map((play, index) => (
          <PlayerCard index={index} play={play} showCards={props.showCards} />
        ))}
      </div>
    </div>
  );
}
