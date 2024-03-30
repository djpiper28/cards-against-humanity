import { useNavigate, useSearchParams } from "@solidjs/router";
import { joinGameUrl } from "../routes";
import { gameIdParamCookie } from "../gameState/gameState";

export default function GameJoinErrorPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const gameId = searchParams[gameIdParamCookie];

  return (
    <div class="flex flex-col justify-center items-center text-2xl">
      <h1>Error joining the game.</h1>
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
