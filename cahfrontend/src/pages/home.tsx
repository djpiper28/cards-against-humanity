import { useNavigate } from "@solidjs/router";
import Card from "../components/gameItems/Card";
import { For, createSignal, onMount } from "solid-js";
import { GameLogicGameInfo } from "../api";
import { apiClient } from "../apiClient";

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
      <div class="flex w-full justify-center">
        <h1 class="text-4xl lg:text-5xl font-bold">Cards Against Humanity</h1>
      </div>
      <button class="w-min" onClick={() => navigate("/create")}>
        <Card
          isWhite={false}
          cardText="Create a Game"
          packName="Click me to make a game"
        />
      </button>
      <div class="flex flex-row flex-wrap gap-5">
        <For each={currentGames()}>
          {(game, _) => (
            <button
              onClick={() => navigate(`/join/${encodeURIComponent(game.id)}`)}
            >
              <Card
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
