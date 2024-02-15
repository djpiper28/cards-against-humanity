export interface Props {
  isWhite: boolean;
  cardText: string;
  packName: string;
}

export default function Card({ isWhite, cardText, packName }: Readonly<Props>) {
  const textColour = isWhite ? "text-black" : "text-white";
  const backgroundColour = isWhite ? "bg-white" : "bg-black";

  return (
    <div
      class={`${textColour} ${backgroundColour} rounded-2xl p-3 md:p6 flex flex-col justify-between h-80 w-64 border-2 text-left`}
      aria-label={`${isWhite ? "white" : "black"} card`}
      data-testid={`${isWhite ? "white" : "black"}`}
    >
      <p class="font-bold text-lg">{cardText}</p>
      <div class="text-xs">@{packName}</div>
    </div>
  );
}
