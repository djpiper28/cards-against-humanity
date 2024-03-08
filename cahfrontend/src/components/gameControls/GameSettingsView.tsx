import { GameSettings } from "../../gameLogicTypes";

interface Props {
  settings: GameSettings;
}

export default function GameSettingsView(props: Readonly<Props>) {
  return (
    <div>
      <h2>Game Settings</h2>
      <div>
        <p>Game Password: {props.settings.gamePassword}</p>
        <p>Max Players: {props.settings.maxPlayers}</p>
        <p>Points To Play To: {props.settings.playingToPoints}</p>
        <p>Max Rounds: {props.settings.maxRounds}</p>
      </div>
      <h3>Card Packs</h3>
      <p>
        {props.settings.cardPacks
          .map((pack) => <span class="font-bold">{pack?.name}</span>)
          .reduce((a, b) => (
            <>
              {a}, {b}
            </>
          ))}
      </p>
    </div>
  );
}
