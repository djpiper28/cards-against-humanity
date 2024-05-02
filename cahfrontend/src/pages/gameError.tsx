import { useNavigate, useSearchParams } from "@solidjs/router";
import { joinGameUrl } from "../routes";
import Header from "../components/typography/Header";
import Button from "../components/buttons/Button";

export const errorMessageParam = "error";

export default function GameJoinErrorPage() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const errorMessage = searchParams[errorMessageParam];

  return (
    <div class="flex flex-col justify-center items-center text-2xl">
      <Header text={errorMessage ?? "Error joining the game"} />
      <Button
        onClick={() => {
          navigate(`${joinGameUrl}`);
        }}
      >
        Try Again
      </Button>
    </div>
  );
}
