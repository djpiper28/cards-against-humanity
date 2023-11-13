import { Api, GameLogicCardPack } from "../api";
import { For, createSignal, onMount } from "solid-js";
import Checkbox from "../components/inputs/Checkbox";

interface Checked {
  checked: boolean;
}

type CardPack = GameLogicCardPack & Checked;

export default function Create() {
  const [packs, setPacks] = createSignal<CardPack[]>([]);
  onMount(async () => {
    const api = new Api();
    const packs = await api.res.packsList();
    const cardPacksList: CardPack[] = [];
    const packData = packs.data;
    for (let cardId in packData) {
      cardPacksList.push({ ...packData[cardId], checked: false });
    }
    setPacks(cardPacksList.sort((a, b) => a.name.localeCompare(b.name)));
  });

  const [selectedPacks, setSelectedPacks] = createSignal<string[]>([]);
  return (
    <>
      <h1>Create A Game of Cards Against Humanity</h1>
      <div class="flex flex-row flex-wrap gap-2 md:gap-1">
        <For each={packs()}>
          {(pack, index) => {
            return (
              <Checkbox
                checked={pack.checked}
                label={`${pack.name} (${
                  pack.whiteCards + pack.blackCards
                } Cards)`}
                onSetChecked={(checked) => {
                  if (checked && !selectedPacks().includes(pack.id)) {
                    setSelectedPacks([...selectedPacks(), pack.id]);
                  } else {
                    setSelectedPacks(
                      selectedPacks().filter((x) => x !== pack.id),
                    );
                  }

                  const newPacks = structuredClone(packs());
                  newPacks[index()].checked = !pack.checked;
                  setPacks(newPacks);
                }}
              />
            );
          }}
        </For>
      </div>
    </>
  );
}
