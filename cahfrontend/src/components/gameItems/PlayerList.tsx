import { GamePlayerList } from "../../gameState/playersList";
import RoundedWhite from "../containers/RoundedWhite";

interface Props {
  players: GamePlayerList;
}

export default function PlayerList(props: Readonly<Props>) {
  return (
    <RoundedWhite>
      {props.players.map((player) => (
        <div class="flex flex-col gap-1">
          <p class="font-bold">{player.name}</p>
          <p class="font-mono">
            {player.connected ? "CONNECTED" : "DISCONNECTED"}
          </p>
        </div>
      ))}
    </RoundedWhite>
  );
}
