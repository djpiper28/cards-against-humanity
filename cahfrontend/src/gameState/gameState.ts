import { GameSettings, GameStateInfo } from "../gameLogicTypes";
import { MsgOnJoin, RpcMessageBody, RpcOnJoinMsg } from "../rpcTypes";
import { WebSocketClient, toWebSocketClient } from "./websocketClient";
import { wsBaseUrl } from "../apiClient";
import WebSocket from "isomorphic-ws";

export const playerIdCookie = "playerId";
export const gameIdParam = "gameId";

class GameState {
  constructor() {}
  private gameId: string = "";
  private playerId: string = "";
  private setup: boolean = false;
  private gameSettings?: GameSettings = undefined;
  private wsClient: WebSocketClient;
  private state: GameStateInfo;

  public setupState(gameId: string, playerId: string) {
    this.gameId = gameId;
    this.playerId = playerId;

    /**
     * Mock new WebSocket when testing
     */
    const ws: WebSocket = new WebSocket(wsBaseUrl);
    this.wsClient = toWebSocketClient(ws, {
      onDisconnect: () => {
        console.error("Disconnected from the game server");
      },
      onConnect: () => {
        console.log("Connected to the game server");
      },
      onReceive: (msg: string) => {
        this.handleRpcMessage(msg);
      },
    });

    this.setup = true;
  }

  public validate(): boolean {
    if (!this.setup) return false;
    if (!this.gameId) return false;
    if (!this.playerId) return false;
    return true;
  }

  private handleOnJoin(msg: RpcOnJoinMsg) {
    this.state = msg.state as GameStateInfo;
  }

  private handleRpcMessage(msg: string): void {
    const rpcMessage = JSON.parse(msg) as RpcMessageBody;
    switch (rpcMessage.type) {
      case MsgOnJoin:
        return this.handleOnJoin(rpcMessage.data as RpcOnJoinMsg);
      default:
        throw new Error(
          `Cannot handle RPC message as type is not valid ${rpcMessage.type}`,
        );
    }
  }
}

export const gameState: GameState = new GameState();
