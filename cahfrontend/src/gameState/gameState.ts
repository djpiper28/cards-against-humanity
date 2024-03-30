import { GameStateInfo } from "../gameLogicTypes";
import { MsgOnJoin, RpcMessageBody, RpcOnJoinMsg } from "../rpcTypes";
import { WebSocketClient, toWebSocketClient } from "./websocketClient";
import { wsBaseUrl } from "../apiClient";
import WebSocket from "isomorphic-ws";

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
  private wsClient: WebSocketClient;
  private state?: GameStateInfo;

  // Events
  public onStateChange?: (state?: GameStateInfo) => void;

  // Logic lmao
  constructor() {}

  public setupState(gameId: string, playerId: string, password: string) {
    this.gameId = gameId;
    this.playerId = playerId;
    this.password = password;

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
  }

  public isOwner(): boolean {
    return this.playerId === this.ownerId;
  }

  private handleOnJoin(msg: RpcOnJoinMsg) {
    this.setState(msg.state as GameStateInfo);
  }

  private handleRpcMessage(msg: string): void {
    console.log(`Received a message ${msg}`);
    const rpcMessage = JSON.parse(msg) as RpcMessageBody;
    switch (rpcMessage.type) {
      case MsgOnJoin:
        console.log("Handling on join message");
        return this.handleOnJoin(rpcMessage.data as RpcOnJoinMsg);
      default:
        throw new Error(
          `Cannot handle RPC message as type is not valid ${rpcMessage.type}`,
        );
    }
  }
}

export const gameState: GameState = new GameState();
