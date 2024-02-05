import { GameSettings } from "../gameLogicTypes";

export const playerIdCookie = "playerId";
export const gameIdParam = "gameID";

class GameState {
  constructor() {}
  private gameId: string = "";
  private playerId: string = "";
  private gameSettings?: GameSettings = undefined;
  private setup: boolean = false;

  public setupState(gameId: string, playerId: string) {
    this.gameId = gameId;
    this.playerId = playerId;
    this.setup = true;
  }

  public validate(): boolean {
    if (!this.setup) return false;
    if (!this.gameId) return false;
    if (!this.playerId) return false;
    // if (!this.gameSettings) return false;
    return true;
  }
}

export const gameState: GameState = new GameState();
