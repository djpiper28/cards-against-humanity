import { For } from "solid-js";
import { GameLogicCardPack } from "../../api";
import Checkbox from "../inputs/Checkbox";

type CardPack = GameLogicCardPack;

interface Props {
  cards: CardPack[];
  selectedPackIds: string[];
  setSelectedPackIds: (packs: string[]) => void;
}

export default function CardsSelector(props: Readonly<Props>) {
  const whiteCards = (): number => {
    return props.selectedPackIds
      .map((x) => props.cards.find((y) => y.id === x))
      .filter((x) => !!x)
      .map((x) => x?.whiteCards ?? 0)
      .reduce((a, b) => a + b, 0);
  };
  const blackCards = (): number => {
    return props.selectedPackIds
      .map((x) => props.cards.find((y) => y.id === x))
      .filter((x) => !!x)
      .map((x) => x?.blackCards ?? 0)
      .reduce((a, b) => a + b, 0);
  };
  const panelTitleCss = () =>
    `text-xl ${
      blackCards() + whiteCards() === 0
        ? "text-error-colour font-bold"
        : "text-black"
    }`;

  return (
    <div class="flex flex-col gap-3">
      <fieldset class="flex flex-row flex-wrap gap-2 md:gap-1 overflow-auto max-h-64">
        <legend class="hidden">Card Packs</legend>
        <For each={props.cards}>
          {(pack) => (
            <Checkbox
              checked={!!props.selectedPackIds.find((x) => x === pack.id)}
              label={`${pack.name} (${
                pack ? (pack.whiteCards ?? 0) + (pack.blackCards ?? 0) : 0
              } Cards)`}
              onSetChecked={(checked) => {
                if (checked && !props.selectedPackIds.includes(pack.id)) {
                  props.setSelectedPackIds([...props.selectedPackIds, pack.id]);
                } else {
                  props.setSelectedPackIds(
                    props.selectedPackIds.filter((x) => x !== pack.id),
                  );
                }
              }}
              secondary={/^CAH:?.*$/i.test(pack.name ?? "")}
            />
          )}
        </For>
      </fieldset>
      <p class={panelTitleCss()}>
        {`You have added ${whiteCards()} white cards and ${blackCards()} black cards.`}
      </p>
    </div>
  );
}
