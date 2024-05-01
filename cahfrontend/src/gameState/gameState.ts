import { GameStateInfo, Player } from "../gameLogicTypes";
import {
  MsgChangeSettings,
  MsgCommandError,
  MsgNewOwner,
  MsgOnJoin,
  MsgOnPlayerCreate,
  MsgOnPlayerDisconnect,
  MsgOnPlayerJoin,
  MsgOnPlayerLeave,
  RpcChangeSettingsMsg,
  RpcCommandErrorMsg,
  RpcMessage,
  RpcMessageBody,
  RpcMessageType,
  RpcNewOwnerMsg,
  RpcOnJoinMsg,
  RpcOnPlayerCreateMsg,
  RpcOnPlayerDisconnectMsg,
  RpcOnPlayerJoinMsg,
  RpcOnPlayerLeaveMsg,
} from "../rpcTypes";
import { WebSocketClient, toWebSocketClient } from "./websocketClient";
import { apiClient, wsBaseUrl } from "../apiClient";
import WebSocket from "isomorphic-ws";
import { GamePlayerList } from "./gamePlayersList";

export const playerIdCookie = "playerId";
/**
 * Used as a cookie name and a search param name. Required for authentication to a game.
 */
export const gameIdParamCookie = "gameId";
export const gamePasswordCookie = "password";
export const authenticationCookie = "authentication";

class GameState {
  private gameId: string = "";
  private playerId: string = "";
  private ownerId: string = "";
  private password: string = "";
  private setup: boolean = false;
  private wsClient?: WebSocketClient;
  private state?: GameStateInfo;
  private players: GamePlayerList = [];

  // Events
  public onStateChange?: (state?: GameStateInfo) => void;
  public onPlayerListChange?: (players: GamePlayerList) => void;
  public onCommandError?: (error: RpcCommandErrorMsg) => void;
  public onChangeSettings?: (settings: RpcChangeSettingsMsg) => void;

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

  public getPlayerId(): string {
    return this.playerId;
  }

  private setState(state: GameStateInfo) {
    this.state = state;
    this.ownerId = state.gameOwnerId;
    this.players = state.players.map((x) => ({
      id: x.id,
      name: x.name,
      connected: true,
      points: x.points,
    }));
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

  public playerList(): GamePlayerList {
    return [...this.players];
  }

  private handleOnPlayerJoin(msg: RpcOnPlayerJoinMsg) {
    const player = this.players.find((x: Player) => x.id === msg.id);
    this.players = this.players.filter((x: Player) => x.id !== msg.id);
    this.players.push({
      id: msg.id,
      name: msg.name,
      connected: true,
      points: player?.points ?? 0,
    });

    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnPlayerCreate(msg: RpcOnPlayerCreateMsg) {
    const player = this.players.find((x: Player) => x.id === msg.id);
    this.players = this.players.filter((x: Player) => x.id !== msg.id);
    this.players.push({
      id: msg.id,
      name: msg.name,
      connected: false,
      points: player?.points ?? 0,
    });

    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnPlayerDisconnect(msg: RpcOnPlayerDisconnectMsg) {
    const player = this.players.find((x: Player) => x.id === msg.id);
    const oldPlayer = this.players.find((x: Player) => x.id === msg.id);
    this.players = this.players.filter((x: Player) => x.id !== msg.id);
    this.players.push({
      id: msg.id,
      name: oldPlayer?.name ?? "error",
      connected: false,
      points: player?.points ?? 0,
    });

    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnPlayerLeave(msg: RpcOnPlayerLeaveMsg) {
    this.players = this.players.filter((x: Player) => x.id !== msg.id);
    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnOwnerChange(msg: RpcNewOwnerMsg) {
    var state: GameStateInfo = structuredClone(
      this.state ?? {
        id: "",
        gameOwnerId: "",
        players: [],
        currentRound: 0,
        creationTime: "",
        gameState: 0,
        currentCardCzarId: "",
        settings: {
          maxRounds: 0,
          maxPlayers: 0,
          playingToPoints: 0,
          gamePassword: "",
          cardPacks: [],
        },
      },
    );
    state.gameOwnerId = msg.id;
    this.setState(state);
  }

  /**
   * Handles an RPC message from the server. When testing call the private method and ignore the "error".
   */
  private handleRpcMessage(msg: string): void {
    const rpcMessage = JSON.parse(msg) as RpcMessageBody;
    switch (rpcMessage.type) {
      case MsgOnJoin:
        console.log("Handling on join message");
        return this.handleOnJoin(rpcMessage.data as RpcOnJoinMsg);
      case MsgOnPlayerJoin:
        console.log("Handling on player join message");
        return this.handleOnPlayerJoin(rpcMessage.data as RpcOnPlayerJoinMsg);
      case MsgOnPlayerCreate:
        console.log("Handling on player create message");
        return this.handleOnPlayerCreate(
          rpcMessage.data as RpcOnPlayerCreateMsg,
        );
      case MsgOnPlayerDisconnect:
        console.log("Handling on player disconnect message");
        return this.handleOnPlayerDisconnect(
          rpcMessage.data as RpcOnPlayerDisconnectMsg,
        );
      case MsgCommandError:
        console.log("Handling command error message");
        return this.onCommandError?.(rpcMessage.data as RpcCommandErrorMsg);
      case MsgChangeSettings:
        console.log("Handling change settings message");
        const data = rpcMessage.data as RpcChangeSettingsMsg;

        if (this.state) {
          const newState: GameStateInfo = this.state;
          newState.settings.maxRounds = data.settings.maxRounds;
          newState.settings.maxPlayers = data.settings.maxPlayers;
          newState.settings.playingToPoints = data.settings.playingToPoints;
          newState.settings.gamePassword = data.settings.gamePassword;
          newState.settings.cardPacks = data.settings.cardPacks;
          this.setState(newState);
        }

        return this.onChangeSettings?.(data);
      case MsgOnPlayerLeave:
        console.log("Handling on player leave message");
        return this.handleOnPlayerLeave(rpcMessage.data as RpcOnPlayerLeaveMsg);
      case MsgNewOwner:
        console.log("Handling new owner message");
        return this.handleOnOwnerChange(rpcMessage.data as RpcNewOwnerMsg);
      default:
        throw new Error(
          `Cannot handle RPC message as type is not valid ${rpcMessage.type}`,
        );
    }
  }

  public encodeMessage(type: RpcMessageType, data: RpcMessage): RpcMessageBody {
    return {
      type: type,
      data: data,
    };
  }

  public sendRpcMessage(type: RpcMessageType, data: RpcMessage) {
    if (!this.wsClient) {
      throw new Error("Cannot send message as websocket is not connected");
    }

    this.onCommandError?.({
      reason: "",
    });
    this.wsClient.sendMessage(JSON.stringify(this.encodeMessage(type, data)));
  }

  public async leaveGame() {
    console.log("Leaving game: ", this.gameId);
    if (!this.wsClient) {
      throw new Error("Cannot leave game as websocket is not connected");
    }

    this.wsClient.disconnect();
    return apiClient.games.leaveDelete();
  }
}

export const gameState: GameState = new GameState();
