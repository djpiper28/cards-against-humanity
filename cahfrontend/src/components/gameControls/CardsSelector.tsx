import { createSignal, For } from "solid-js";
import { GameLogicCardPack } from "../../api";
import Checkbox from "../inputs/Checkbox";
import Input, { InputType } from "../inputs/Input";

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
    `text - xl ${
      blackCards() + whiteCards() === 0
        ? "text-error-colour font-bold"
        : "text-black"
    }
  `;
  const [query, setQuery] = createSignal("");

  const normaliseString = (str: string): string =>
    str.toLowerCase().replace(/\s{2,}/, " ");
  const filterByQuery = (
    cardPacks: GameLogicCardPack[],
  ): GameLogicCardPack[] => {
    const normalisedQuery = normaliseString(query());
    return cardPacks.filter((card) =>
      normaliseString(card.name).includes(normalisedQuery),
    );
  };

  return (
    <div class="flex flex-col gap-3 max-h-72 lg:max-h-96 2xl:max-h-screen">
      <Input
        inputType={InputType.Text}
        label="Search"
        placeHolder="Search..."
        onChanged={setQuery}
        value={query()}
      />
      <fieldset class="flex flex-row flex-wrap gap-2 md:gap-1 overflow-auto">
        <legend class="hidden">Card Packs</legend>
        <For each={filterByQuery(props.cards)}>
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
