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
      <h1 class="text-2xl font-bold">
        Create A Game of Cards Against Humanity
      </h1>
      <h2 class="text-xl">Choose Some Card Packs</h2>
      <div class="flex flex-col gap-5 rounded-2xl border-2 p-5 bg-gray-100">
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

        <p class="text-xl">
          {`You have ${selectedPacks()
            .map((x) => packs().find((y) => y.id === x))
            .filter((x) => !!x)
            .map((x) => x.whiteCards)
            .reduce((a, b) => a + b, 0)} white cards and ${selectedPacks()
            .map((x) => packs().find((y) => y.id === x))
            .filter((x) => !!x)
            .map((x) => x.blackCards)
            .reduce((a, b) => a + b, 0)} black cards.`}
        </p>
      </div>
    </>
  );
}
