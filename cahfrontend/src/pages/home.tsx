import { useNavigate } from "@solidjs/router";
import Card from "../components/gameItems/Card";

export default function Home() {
  const navigate = useNavigate();
  return (
    <>
      <div class="flex w-full justify-center">
        <h1 class="text-4xl lg:text-5xl font-bold">Cards Against Humanity</h1>
      </div>
      <button onClick={() => navigate("/create")}>
        <Card
          isWhite={false}
          cardText="Create a Game"
          packName="Click me to make a game"
        />
      </button>
    </>
  );
}
