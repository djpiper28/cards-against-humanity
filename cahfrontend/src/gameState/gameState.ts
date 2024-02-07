import { GameSettings } from "../gameLogicTypes";
import { RpcMessageBody } from "../rpcTypes";
import { WebSocketClient, toWebSocketClient } from "./websocketClient";
import { wsBaseUrl } from "../apiClient";
import WebSocket from "isomorphic-ws";

export const playerIdCookie = "playerId";
export const gameIdParam = "gameID";

class GameState {
  constructor() {}
  private gameId: string = "";
  private playerId: string = "";
  private setup: boolean = false;
  private gameSettings?: GameSettings = undefined;
  private wsClient: WebSocketClient;

  public setupState(gameId: string, playerId: string) {
    this.gameId = gameId;
    this.playerId = playerId;

    /**
     * Mock new WebSocket when testing
     */
    const ws: WebSocket = new WebSocket(wsBaseUrl);
    this.wsClient = toWebSocketClient(ws, {
      onDisconnect: () => {},
      onConnect: () => {},
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

  private handleRpcMessage(msg: string): void {
    const rpcMessage = JSON.parse(msg) as RpcMessageBody;
    switch (rpcMessage.type) {
      default:
        throw new Error(
          `Cannot handle RPC message as type is not valid ${rpcMessage.type}`,
        );
    }
  }
}

export const gameState: GameState = new GameState();
