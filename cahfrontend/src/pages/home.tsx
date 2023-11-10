import Card from "../components/gameItems/Card";

export default function Home() {
  return (
    <>
      <div class="flex w-full justify-center">
        <h1 class="text-4xl lg:text-5xl font-bold">Cards Against Humanity</h1>
      </div>
      <button>
        <Card
          isWhite={true}
          cardText="Create a Game"
          packName="Click me to make a game"
        />
      </button>
    </>
  );
}
