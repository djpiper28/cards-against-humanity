import { Api, GameLogicCardPack } from "../api";
import { For, createSignal, onMount } from "solid-js";
import Checkbox from "../components/inputs/Checkbox";
import Input, { InputType } from "../components/inputs/Input";
import { useNavigate } from "@solidjs/router";
import {
  MaxPasswordLength,
  MaxPlayerNameLength,
  MinPlayerNameLength,
  MinPlayers,
  MaxPlayers,
  MinPlayingToPoints,
  MaxPlayingToPoints,
  MinRounds,
  MaxRounds,
} from "../gameLogicTypes";
import { gameIdParam, playerIdCookie } from "../gameState/gameState";
import { cookieStorage } from "@solid-primitives/storage";

interface Checked {
  checked: boolean;
}

type CardPack = GameLogicCardPack & Checked;

export default function Create() {
  const navigate = useNavigate();
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
  const [gamePassword, setGamePassword] = createSignal("");
  const [maxPlayers, setMaxPlayers] = createSignal(6);
  const [gameRounds, setGameRounds] = createSignal(25);
  const [playingToPoints, setPlayingToPoints] = createSignal(10);
  const [playerName, setPlayerName] = createSignal("");
  const [errorMessage, setErrorMessage] = createSignal("");

  const editPanelCss =
    "flex flex-col gap-5 rounded-2xl border-2 p-3 md:p-5 bg-gray-100";
  const panelTitleCss = "text-xl";
  return (
    <>
      <h1 class="text-2xl font-bold">
        Create A Game of Cards Against Humanity
      </h1>
      <div class={editPanelCss}>
        <h2
          class={`${panelTitleCss} ${
            selectedPacks().length === 0 ? "text-error-colour" : ""
          }`}
        >
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
                  secondary={/^CAH:.*$/i.test(pack.name ?? "")}
                />
              );
            }}
          </For>
        </div>

        <p class={panelTitleCss}>
          {`You have added ${selectedPacks()
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

      <div class={editPanelCss}>
        <h2 class={panelTitleCss}>Other Game Settings</h2>
        <div class="flex flex-row flex-wrap gap-2 md:gap-1">
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

          <Input
            inputType={InputType.Text}
            placeHolder="password"
            value={gamePassword()}
            onChanged={setGamePassword}
            label="Game Password"
            errorState={gamePassword().length > MaxPasswordLength}
          />

          <Input
            inputType={InputType.PositiveNumber}
            placeHolder="max players"
            value={maxPlayers().toString()}
            onChanged={(text) => setMaxPlayers(parseInt(text))}
            label="Max Players"
            errorState={maxPlayers() < MinPlayers || maxPlayers() > MaxPlayers}
          />

          <Input
            inputType={InputType.PositiveNumber}
            placeHolder="points to play to"
            value={playingToPoints().toString()}
            onChanged={(text) => setPlayingToPoints(parseInt(text))}
            label="Points To Play To"
            errorState={
              playingToPoints() < MinPlayingToPoints ||
              playingToPoints() > MaxPlayingToPoints
            }
          />

          <Input
            inputType={InputType.PositiveNumber}
            placeHolder="max game rounds"
            value={gameRounds().toString()}
            onChanged={(text) => setGameRounds(parseInt(text))}
            label="Max Game Rounds"
            errorState={gameRounds() < MinRounds || gameRounds() > MaxRounds}
          />
        </div>

        <p class="text-error-colour text-lg">{errorMessage()}</p>

        <button
          onclick={() => {
            const api = new Api();
            api.games
              .createCreate({
                settings: {
                  cardPacks: selectedPacks(),
                  gamePassword: gamePassword(),
                  maxPlayers: maxPlayers(),
                  maxRounds: gameRounds(),
                  playingToPoints: playingToPoints(),
                },
                playerName: playerName(),
              })
              .then((newGame) => {
                console.log("Creating game for ", JSON.stringify(newGame.data));
                cookieStorage.setItem(playerIdCookie, newGame.data.playerId);
                navigate(
                  `/join?${gameIdParam}=${encodeURIComponent(
                    newGame.data.gameId,
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
