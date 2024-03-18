import { GameStateInfo } from "../gameLogicTypes";
import { MsgOnJoin, RpcMessageBody, RpcOnJoinMsg } from "../rpcTypes";
import { WebSocketClient, toWebSocketClient } from "./websocketClient";
import { wsBaseUrl } from "../apiClient";
import WebSocket from "isomorphic-ws";

export const playerIdCookie = "playerId";
export const gameIdParam = "gameId";

class GameState {
  private gameId: string = "";
  private playerId: string = "";
  private setup: boolean = false;
  private wsClient: WebSocketClient;
  private state?: GameStateInfo;

  // Events
  public onStateChange?: (state?: GameStateInfo) => void;

  // Logic lmao
  constructor() {}

  public setupState(gameId: string, playerId: string) {
    this.gameId = gameId;
    this.playerId = playerId;

    const url = `${wsBaseUrl}?game_id=${encodeURIComponent(
      gameId
    )}&player_id=${encodeURIComponent(playerId)}`;
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
    this.emitState();
  }

  public emitState() {
    this.onStateChange?.(structuredClone(this.state));
  }

  private handleOnJoin(msg: RpcOnJoinMsg) {
    this.setState(msg.state as GameStateInfo);
  }

  private handleRpcMessage(msg: string): void {
    console.log(`Received a message ${msg}`);
    const rpcMessage = JSON.parse(msg) as RpcMessageBody;
    switch (rpcMessage.type) {
      case MsgOnJoin:
        console.log("Handling on join message")
        return this.handleOnJoin(rpcMessage.data as RpcOnJoinMsg);
      default:
        throw new Error(
          `Cannot handle RPC message as type is not valid ${rpcMessage.type}`
        );
    }
  }
}

export const gameState: GameState = new GameState();
