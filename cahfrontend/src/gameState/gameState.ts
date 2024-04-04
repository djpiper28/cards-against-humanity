import { GameStateInfo } from "../gameLogicTypes";
import {
  MsgOnJoin,
  MsgOnPlayerJoin,
  RpcMessageBody,
  RpcOnJoinMsg,
  RpcOnPlayerJoinMsg,
} from "../rpcTypes";
import { WebSocketClient, toWebSocketClient } from "./websocketClient";
import { wsBaseUrl } from "../apiClient";
import WebSocket from "isomorphic-ws";
import { PlayerList } from "./playersList";

export const playerIdCookie = "playerId";
/**
 * Used as a cookie name and a search param name. Required for authentication to a game.
 */
export const gameIdParamCookie = "gameId";
export const gamePasswordCookie = "password";

class GameState {
  private gameId: string = "";
  private playerId: string = "";
  private ownerId: string = "";
  private password: string = "";
  private setup: boolean = false;
  private wsClient?: WebSocketClient;
  private state?: GameStateInfo;
  private players: PlayerList = [];

  // Events
  public onStateChange?: (state?: GameStateInfo) => void;
  public onPlayerListChange?: (players: PlayerList) => void;

  // Logic lmao
  constructor() {}

  public setupState(gameId: string, playerId: string, password: string) {
    this.gameId = gameId;
    this.playerId = playerId;
    this.password = password;
    this.players = [];

    this.onStateChange = undefined;
    this.onPlayerListChange = undefined;

    const url = wsBaseUrl;
    console.log(`Connecting to ${url}`);

    /**
     * Mock new WebSocket when testing
     */
    const ws: WebSocket = new WebSocket(url);
    this.wsClient = toWebSocketClient(ws, {
      onDisconnect: () => {
        console.error("Disconnected from the game server");
      },
      onConnect: () => {
        console.log("Connected to the game server");
      },
      onReceive: (msg: string) => {
        console.log(`Received a message ${msg}`);
        this.handleRpcMessage(msg);
      },
    });

    console.log("State Setup");
    this.setup = true;
  }

  public validate(): boolean {
    if (!this.setup) return false;
    if (!this.gameId) return false;
    if (!this.playerId) return false;
    return true;
  }

  private setState(state: GameStateInfo) {
    this.state = state;
    this.ownerId = state.gameOwnerId;
    this.emitState();
  }

  public emitState() {
    this.onStateChange?.(structuredClone(this.state));
    this.onPlayerListChange?.(this.playerList());
  }

  public isOwner(): boolean {
    return this.playerId === this.ownerId;
  }

  private handleOnJoin(msg: RpcOnJoinMsg) {
    this.setState(msg.state as GameStateInfo);
  }

  public playerList(): PlayerList {
    return structuredClone(this.players);
  }

  private handleOnPlayerJoin(msg: RpcOnPlayerJoinMsg) {
    this.players = this.players.filter((x) => x.id !== msg.id);
    this.players.push({
      id: msg.id,
      name: msg.name,
      connected: true,
    });

    this.onPlayerListChange?.(this.playerList());
  }

  /**
   * Handles an RPC message from the server. When testing call the private method and ignore the "error".
   */
  private handleRpcMessage(msg: string): void {
    console.log(`Received a message ${msg}`);
    const rpcMessage = JSON.parse(msg) as RpcMessageBody;
    switch (rpcMessage.type) {
      case MsgOnJoin:
        console.log("Handling on join message");
        return this.handleOnJoin(rpcMessage.data as RpcOnJoinMsg);
      case MsgOnPlayerJoin:
        console.log("Handling on player join message");
        return this.handleOnPlayerJoin(rpcMessage.data as RpcOnPlayerJoinMsg);
      default:
        throw new Error(
          `Cannot handle RPC message as type is not valid ${rpcMessage.type}`,
        );
    }
  }
}

export const gameState: GameState = new GameState();
