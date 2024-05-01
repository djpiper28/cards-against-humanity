import { GameLogicCardPack } from "../api";
import { createSignal, onMount } from "solid-js";
import Input, { InputType } from "../components/inputs/Input";
import { useNavigate } from "@solidjs/router";
import { MaxPlayerNameLength, MinPlayerNameLength } from "../gameLogicTypes";
import {
  authenticationCookie,
  gameIdParamCookie,
  gamePasswordCookie,
  gameState,
  playerIdCookie,
} from "../gameState/gameState";
import { cookieStorage } from "@solid-primitives/storage";
import { apiClient, cookieOptions } from "../apiClient";
import GameSettingsInput, {
  Settings,
} from "../components/gameControls/GameSettingsInput";
import { validate as validateGameSettings } from "../components/gameControls/GameSettingsInputValidation";
import { joinGameUrl } from "../routes";
import RoundedWhite from "../components/containers/RoundedWhite";
import Header from "../components/typography/Header";
import Button from "../components/buttons/Button";
import CardsSelector from "../components/gameControls/CardsSelector";
import clearGameCookies from "../gameState/clearGameCookies";
import SubHeader from "../components/typography/SubHeader";

export default function Create() {
  const navigate = useNavigate();
  const [packs, setPacks] = createSignal<GameLogicCardPack[]>([]);
  const [errorMessage, setErrorMessage] = createSignal("");
  onMount(async () => {
    try {
      const packs = await apiClient.res.packsList();
      const cardPacksList: GameLogicCardPack[] = [];
      const packData = packs.data;
      for (let cardId in packData) {
        cardPacksList.push({ ...packData[cardId] });
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

  return (
    <>
      <Header text="Create a game" />
      <RoundedWhite>
        <SubHeader
          text={`Choose Some Card Packs ${
            selectedPacks().length === 0 ? "(No Packs Selected)" : ""
          }`}
        />
        <CardsSelector
          setSelectedPackIds={setSelectedPacks}
          selectedPackIds={selectedPacks()}
          cards={packs()}
        />
      </RoundedWhite>

      <RoundedWhite>
        <SubHeader text="Other Game Settings" />
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

        <Button
          onClick={async () => {
            if (!validateGameSettings(settings)) {
              console.error("The game settings are invalid");
              return;
            }

            try {
              await gameState.leaveGame();
              clearGameCookies();
            } catch (e) {
              console.log(e);
            }

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
                  cookieOptions,
                );
                cookieStorage.setItem(
                  gameIdParamCookie,
                  newGame.data.gameId ?? "error",
                  cookieOptions,
                );
                cookieStorage.setItem(
                  gamePasswordCookie,
                  gameSettings().gamePassword,
                  cookieOptions,
                );
                cookieStorage.setItem(
                  authenticationCookie,
                  newGame.data.authentication,
                  cookieOptions,
                );
                navigate(
                  `${joinGameUrl}?${gameIdParamCookie}=${encodeURIComponent(
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
        </Button>
      </RoundedWhite>
    </>
  );
}
