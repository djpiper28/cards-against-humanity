import { useNavigate } from "@solidjs/router";
import Card from "../components/gameItems/Card";
import { For, createSignal, onMount } from "solid-js";
import { GameLogicGameInfo } from "../api";
import { apiClient } from "../apiClient";
import { createGameUrl, joinGameUrl } from "../routes";
import { gameIdParamCookie } from "../gameState/gameState";
import Header from "../components/typography/Header";

export default function Home() {
  const navigate = useNavigate();
  const [currentGames, setCurrentGames] = createSignal<GameLogicGameInfo[]>([]);

  onMount(async () => {
    const resp = await apiClient.games.notFullList();
    setCurrentGames(
      resp.data.slice(0, 7).sort((a, b) => {
        if (a.hasPassword) {
          return 1;
        }
        return a.id < b.id ? -1 : 1;
      }),
    );
  });

  return (
    <>
      <button class="w-min" onClick={() => navigate(createGameUrl)}>
        <Card
          id={"create"}
          isWhite={false}
          cardText="Create a Game"
          packName="Click me to make a game"
        />
      </button>
      <Header
        text={
          currentGames().length === 0 ? "No games in progress" : "Join a Game"
        }
      />
      <div class="flex flex-row flex-wrap gap-5">
        <For each={currentGames()}>
          {(game, index) => (
            <button
              onClick={() =>
                navigate(
                  `${joinGameUrl}?${gameIdParamCookie}=${encodeURIComponent(game.id ?? "error")}`,
                )
              }
            >
              <Card
                id={index()}
                isWhite={true}
                cardText={`${game.playerCount}/${
                  game.maxPlayers
                } Players. Playing To ${game.playingTo}. ${
                  game.hasPassword ? "Password Protected" : ""
                }`}
                packName="Click to join game"
              />
            </button>
          )}
        </For>
      </div>
    </>
  );
}
