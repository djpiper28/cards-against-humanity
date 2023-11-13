import { Api, GameLogicCardPack, MainGetPacksResp } from "../api";
import { createEffect, createSignal } from "solid-js";
import Checkbox from "../components/inputs/Checkbox";

export default function Create() {
  const [packs, setPacks] = createSignal<GameLogicCardPack[]>([]);
  createEffect(async () => {
    const api = new Api();
    const packs = await api.res.packsList();
    const cardPacksList: GameLogicCardPack[] = [];
    const packData = packs.data;
    for (let cardId in packData) {
      cardPacksList.push(packData[cardId]);
    }
    setPacks(cardPacksList);
  }, []);

  return (
    <>
      <h1>Create A Game of Cards Against Humanity</h1>
      <div class="flex flex=row flex-wrap">
        {packs()
          .sort((a, b) => a.name.localeCompare(b.name))
          .map((pack: GameLogicCardPack) => (
            <Checkbox
              label={`${pack.name} (${
                pack.whiteCards + pack.blackCards
              } Cards)`}
              checked={false}
              onSetChecked={console.log}
            />
          ))}
      </div>
    </>
  );
}
