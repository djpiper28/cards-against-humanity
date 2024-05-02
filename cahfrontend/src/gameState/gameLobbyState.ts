import { GameSettings, GameState } from "../gameLogicTypes";

export interface GameLobbyState {
  ownerId: string;
  settings: GameSettings;
  creationTime: Date;
  gameState: GameState;
}
