import {
  Settings,
  validateGamePassword,
  validateMaxGameRounds,
  validateMaxPlayers,
  validatePointsToPlayTo,
} from "./GameSettingsInput";

export function validate(settings: Settings): boolean {
  return (
    validateGamePassword(settings.gamePassword) &&
    validateMaxPlayers(settings.maxPlayers) &&
    validatePointsToPlayTo(settings.playingToPoints) &&
    validateMaxGameRounds(settings.maxRounds)
  );
}
