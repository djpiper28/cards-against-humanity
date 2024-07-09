export interface Props {
  isWhite: boolean;
  cardText: string;
  packName: string;
}

export default function Card(props: Readonly<Props>) {
  const textColour = props.isWhite ? "text-black" : "text-white";
  const backgroundColour = props.isWhite ? "bg-white" : "bg-black";

  return (
    <div
      class={`${textColour} ${backgroundColour} rounded-2xl p-3 md:p6 flex flex-col justify-between h-60 w-52 border-2 text-left`}
      aria-label={`${props.isWhite ? "white" : "black"} card`}
      data-testid={`${props.isWhite ? "white" : "black"}`}
    >
      <p class="font-bold text-lg">{props.cardText}</p>
      <div class="text-xs">{props.packName}</div>
    </div>
  );
}
