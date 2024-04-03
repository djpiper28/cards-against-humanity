import { useNavigate, useSearchParams } from "@solidjs/router";
import { joinGameUrl } from "../routes";
import { gameIdParamCookie } from "../gameState/gameState";
import Header from "../components/typography/Header";

export default function GameJoinErrorPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const gameId = searchParams[gameIdParamCookie];

  return (
    <div class="flex flex-col justify-center items-center text-2xl">
      <Header text="Error joining the game" />
      <button
        onclick={() => {
          navigate(`${joinGameUrl}?${gameIdParamCookie}=${gameId}`);
        }}
      >
        Try Again
      </button>
    </div>
  );
}
