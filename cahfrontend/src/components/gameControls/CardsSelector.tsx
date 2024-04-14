import { GameLogicCardPack } from "../../api";
import Checkbox from "../inputs/Checkbox";

type CardPack = GameLogicCardPack;

interface Props {
  cards: CardPack[];
  selectedPackIds: string[];
  setSelectedPackIds: (packs: string[]) => void;
}

export default function CardsSelector(props: Readonly<Props>) {
  return (
    <div class="flex flex-row flex-wrap gap-2 md:gap-1 overflow-auto max-h-64">
      {props.cards.map((pack) => (
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
      ))}
    </div>
  );
}
