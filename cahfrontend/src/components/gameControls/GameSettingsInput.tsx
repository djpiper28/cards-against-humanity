import Input, { InputType } from "../inputs/Input";
import {
  MaxPasswordLength,
  MinPlayers,
  MaxPlayers,
  MinPlayingToPoints,
  MaxPlayingToPoints,
  MinRounds,
  MaxRounds,
  GameSettings,
} from "../../gameLogicTypes";

export type Settings = Omit<GameSettings, "cardPacks">;

interface Props {
  settings: Settings;
  setSettings: (settings: Settings) => void;
  errorMessage: string;
}

export default function GameSettingsInput(props: Readonly<Props>) {
  return (
    <div class="flex flex-col gap-2">
      <div class="flex flex-row flex-wrap gap-2 md:gap-1">
        <Input
          inputType={InputType.Text}
          placeHolder="password"
          value={props.settings.gamePassword}
          onChanged={(password: string) => {
            props.settings.gamePassword = password;
            props.setSettings(props.settings);
          }}
          label="Game Password"
          errorState={!validateGamePassword(props.settings.gamePassword)}
        />

        <Input
          inputType={InputType.PositiveNumber}
          placeHolder="max players"
          value={props.settings.maxPlayers.toString()}
          onChanged={(text: string) => {
            const num = parseInt(text);
            props.settings.maxPlayers = num;
            props.setSettings(props.settings);
          }}
          label="Max Players"
          errorState={!validateMaxPlayers(props.settings.maxPlayers)}
        />

        <Input
          inputType={InputType.PositiveNumber}
          placeHolder="points to play to"
          value={props.settings.playingToPoints.toString()}
          onChanged={(text: string) => {
            const num = parseInt(text);
            props.settings.playingToPoints = num;
            props.setSettings(props.settings);
          }}
          label="Points To Play To"
          errorState={!validatePointsToPlayTo(props.settings.playingToPoints)}
        />

        <Input
          inputType={InputType.PositiveNumber}
          placeHolder="max game rounds"
          value={props.settings.maxRounds.toString()}
          onChanged={(text: string) => {
            const num = parseInt(text);
            props.settings.maxRounds = num;
            props.setSettings(props.settings);
          }}
          label="Max Game Rounds"
          errorState={!validateMaxGameRounds(props.settings.maxRounds)}
        />
      </div>

      <p class="text-error-colour text-lg">{props.errorMessage}</p>
    </div>
  );
}

export function validateGamePassword(password: string): boolean {
  return password.length <= MaxPasswordLength;
}

export function validateMaxPlayers(maxPlayers: number): boolean {
  return maxPlayers >= MinPlayers && maxPlayers <= MaxPlayers;
}

export function validatePointsToPlayTo(points: number): boolean {
  return points >= MinPlayingToPoints && points <= MaxPlayingToPoints;
}

export function validateMaxGameRounds(maxRounds: number): boolean {
  return maxRounds >= MinRounds && maxRounds <= MaxRounds;
}
