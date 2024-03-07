import { GameLogicCardPack } from "../api";
import { For, createSignal, onMount } from "solid-js";
import Checkbox from "../components/inputs/Checkbox";
import Input, { InputType } from "../components/inputs/Input";
import { useNavigate } from "@solidjs/router";
import { MaxPlayerNameLength, MinPlayerNameLength } from "../gameLogicTypes";
import { gameIdParam, playerIdCookie } from "../gameState/gameState";
import { cookieStorage } from "@solid-primitives/storage";
import { apiClient } from "../apiClient";
import GameSettingsInput, {
  Settings,
} from "../components/gameControls/GameSettingsInput";

interface Checked {
  checked: boolean;
}

type CardPack = GameLogicCardPack & Checked;

export default function Create() {
  const navigate = useNavigate();
  const [packs, setPacks] = createSignal<CardPack[]>([]);
  const [errorMessage, setErrorMessage] = createSignal("");
  onMount(async () => {
    try {
      const packs = await apiClient.res.packsList();
      const cardPacksList: CardPack[] = [];
      const packData = packs.data;
      for (let cardId in packData) {
        cardPacksList.push({ ...packData[cardId], checked: false });
      }
      setPacks(
        cardPacksList.sort((a, b) => {
          if (!a.name || !b.name) return 0;
          return a.name.localeCompare(b.name);
        }),
      );
    } catch (err) {
      console.error(err);
      setErrorMessage(`Error getting card packs: ${err}`);
    }
  });

  const settings: Settings = {
    gamePassword: "",
    maxPlayers: 6,
    maxRounds: 25,
    playingToPoints: 10,
  };
  const [selectedPacks, setSelectedPacks] = createSignal<string[]>([]);
  const [playerName, setPlayerName] = createSignal("");
  const [gameSettings, setGameSettings] = createSignal(settings);

  const whiteCards = (): number => {
    return selectedPacks()
      .map((x) => packs().find((y) => y.id === x))
      .filter((x) => !!x)
      .map((x) => x?.whiteCards ?? 0)
      .reduce((a, b) => a + b, 0);
  };
  const blackCards = (): number => {
    return selectedPacks()
      .map((x) => packs().find((y) => y.id === x))
      .filter((x) => !!x)
      .map((x) => x?.whiteCards ?? 0)
      .reduce((a, b) => a + b, 0);
  };

  const editPanelCss =
    "flex flex-col gap-5 rounded-2xl border-2 p-3 md:p-5 bg-gray-100";
  const panelTitleCss = () =>
    `text-xl ${blackCards() + whiteCards() === 0 ? "text-error-colour font-bold" : "text-black"}`;
  return (
    <>
      <h1 class="text-2xl font-bold">Create A Game</h1>
      <div class={editPanelCss}>
        <h2 class={`${panelTitleCss()}`}>
          {`Choose Some Card Packs ${
            selectedPacks().length === 0 ? "(No Packs Selected)" : ""
          }`}
        </h2>
        <div class="flex flex-row flex-wrap gap-2 md:gap-1 overflow-auto max-h-64">
          <For each={packs()}>
            {(pack, index) => {
              return (
                <Checkbox
                  checked={pack.checked}
                  label={`${pack.name} (${
                    pack ? (pack.whiteCards ?? 0) + (pack.blackCards ?? 0) : 0
                  } Cards)`}
                  onSetChecked={(checked) => {
                    if (checked && !selectedPacks().includes(pack.id ?? "")) {
                      setSelectedPacks([...selectedPacks(), pack.id ?? ""]);
                    } else {
                      setSelectedPacks(
                        selectedPacks().filter((x) => x !== pack.id),
                      );
                    }

                    const newPacks = structuredClone(packs());
                    newPacks[index()].checked = !pack.checked;
                    setPacks(newPacks);
                  }}
                  secondary={/^CAH:?.*$/i.test(pack.name ?? "")}
                />
              );
            }}
          </For>
        </div>

        <p class={panelTitleCss()}>
          {`You have added ${whiteCards()} white cards and ${blackCards()} black cards.`}
        </p>
      </div>

      <div class={editPanelCss}>
        <h2 class={panelTitleCss()}>Other Game Settings</h2>
        <Input
          inputType={InputType.Text}
          placeHolder="John Smith"
          value={playerName()}
          onChanged={setPlayerName}
          label="Player Name"
          autocomplete="name"
          errorState={
            playerName().length < MinPlayerNameLength ||
            playerName().length > MaxPlayerNameLength
          }
        />

        <GameSettingsInput
          settings={gameSettings()}
          errorMessage={errorMessage()}
          setSettings={setGameSettings}
        />

        <button
          onclick={() => {
            apiClient.games
              .createCreate({
                settings: {
                  cardPacks: selectedPacks(),
                  gamePassword: gameSettings().gamePassword,
                  maxPlayers: gameSettings().maxPlayers,
                  maxRounds: gameSettings().maxRounds,
                  playingToPoints: gameSettings().playingToPoints,
                },
                playerName: playerName(),
              })
              .then((newGame) => {
                console.log("Creating game for ", JSON.stringify(newGame.data));
                cookieStorage.setItem(
                  playerIdCookie,
                  newGame.data.playerId ?? "error",
                );
                navigate(
                  `/join?${gameIdParam}=${encodeURIComponent(
                    newGame.data.gameId ?? "error",
                  )}`,
                );
              })
              .catch((err) => {
                setErrorMessage(err.error.error);
              });
          }}
        >
          Create Game
        </button>
      </div>
    </>
  );
}
