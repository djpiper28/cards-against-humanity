import { useNavigate } from "@solidjs/router";
import { joinGameUrl } from "../routes";

export default function GameJoinErrorPage() {
  const navigate = useNavigate();
  return (
    <div class="flex flex-col justify-center items-center text-2xl">
      <h1>Error joining the game.</h1>
      <button onclick={() => {
        navigate(joinGameUrl);
      }}>Try Again</button>
    </div>
  );
}
