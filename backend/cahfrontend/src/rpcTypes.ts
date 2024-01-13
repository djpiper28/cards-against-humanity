// Code generated by tygo. DO NOT EDIT.

//////////
// source: connection.go

export interface GameMessage {
  Message: string;
  GameId: string /* uuid */;
  PlayerId: string /* uuid */;
}
export interface WsConnection {
  Conn: NetworkConnection;
  PlayerId: string /* uuid */;
  GameID: string /* uuid */;
  JoinTime: string /* RFC3339 */;
  LastPingTime: string /* RFC3339 */;
  WsRecieve: any;
  WsBroadcast: any;
}

//////////
// source: connection_manager.go

export interface GameConnection {
  PlayerConnectionMap: { [key: string /* uuid */]: WsConnection | undefined};
}
export interface GlobalConnectionManager {
  /**
   * Maps a game ID to the game connection pool
   */
  GameConnectionMap: { [key: string /* uuid */]: GameConnection | undefined};
}

//////////
// source: connection_manager.test.go

export interface MockConnection {
}

//////////
// source: game_metrics.go

export interface Metrics {
}

//////////
// source: games_endpoints.go

export interface GameCreateSettings {
  maxRounds: number /* uint */;
  playingToPoints: number /* uint */;
  gamePassword: string;
  maxPlayers: number /* uint */;
  cardPacks: string /* uuid */[];
}
export interface GameCreateRequest {
  playerName: string;
  settings: GameCreateSettings;
}
export interface GameCreatedResp {
  gameId: string /* uuid */;
  playerId: string /* uuid */;
}
export const JoinGameGameIdParam = "game_id";
export const JoinGamePlayerIdParam = "player_id";

//////////
// source: helpers.go

export interface ApiError {
  error: string;
}

//////////
// source: network_connection.go

/**
 * An interface to abstract a connection over a network to allow for networking to
 * be tested with mocked connections.
 */
export type NetworkConnection = any;
export interface WebsocketConnection {
  Conn?: any /* websocket.Conn */;
}

//////////
// source: resources_endpoints.go


//////////
// source: rpc.go

export type RpcMessageType = number /* int */;
/**
 * The type of the message
 */
export const MsgOnJoin = 0;
export interface RpcMessageBody {
  type: RpcMessageType;
}
export type RpcMessage = any;
export interface RpcOnJoinMsg {
  whiteCardCount: number /* int */;
  blackCardCount: number /* int */;
  playerNames: string[];
  cardPacks: string /* uuid */[];
  GameSettings: any /* gameLogic.GameSettings */;
  blackCard?: any /* gameLogic.BlackCard */;
  WhiteCardsPlayed: number /* int */;
  currentCardCzar: string;
  gameOwner: string;
  gameState: any /* gameLogic.GameState */;
  currentRound: number /* uint */;
}
