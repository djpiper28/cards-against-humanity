/**
 * A player in a game
 */
export interface Player {
  readonly id: string;
  readonly name: string;
  readonly connected: boolean;
  readonly points: number;
}

/**
 * A list of players who are in a game
 */
export type GamePlayerList = Player[];
