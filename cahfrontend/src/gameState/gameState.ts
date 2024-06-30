import {
  GameStateInfo,
  GameStateWhiteCardsBeingSelected,
  Player,
} from "../gameLogicTypes";
import {
  MsgChangeSettings,
  MsgCommandError,
  MsgNewOwner,
  MsgOnJoin,
  MsgOnPlayerCreate,
  MsgOnPlayerDisconnect,
  MsgOnPlayerJoin,
  MsgOnPlayerLeave,
  MsgPing,
  MsgPlayCards,
  MsgRoundInformation,
  MsgStartGame,
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
  RpcPlayCardsMsg,
  RpcRoundInformationMsg,
} from "../rpcTypes";
import { WebSocketClient, toWebSocketClient } from "./websocketClient";
import { apiClient, wsBaseUrl } from "../apiClient";
import WebSocket from "isomorphic-ws";
import { GamePlayerList } from "./gamePlayersList";
import { GameLobbyState } from "./gameLobbyState";

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
  private password: string = "";
  private setup: boolean = false;
  private wsClient?: WebSocketClient;

  private players: GamePlayerList = [];
  private lobbyState: GameLobbyState = {
    ownerId: "",
    settings: {
      gamePassword: "",
      maxPlayers: 0,
      cardPacks: [],
      maxRounds: 0,
      playingToPoints: 0,
    },
    creationTime: new Date(),
    gameState: 0,
  };
  private roundState: RpcRoundInformationMsg = {
    roundNumber: 0,
    currentCardCzarId: "",
    blackCard: {
      id: 0,
      cardsToPlay: 0,
      bodyText: "",
    },
    yourHand: [],
  };

  // Events
  public onLobbyStateChange?: (state?: GameLobbyState) => void;
  public onRoundStateChange?: (state?: RpcRoundInformationMsg) => void;
  public onPlayerListChange?: (players: GamePlayerList) => void;
  public onCommandError?: (error: RpcCommandErrorMsg) => void;
  public onChangeSettings?: (settings: RpcChangeSettingsMsg) => void;
  public onError?: (error: string) => void;

  // Logic lmao
  constructor() {}

  public setupState(gameId: string, playerId: string, password: string) {
    this.gameId = gameId;
    this.playerId = playerId;
    this.password = password;
    this.players = [];

    this.lobbyState = {
      ownerId: "",
      settings: {
        gamePassword: "",
        maxPlayers: 0,
        cardPacks: [],
        maxRounds: 0,
        playingToPoints: 0,
      },
      creationTime: new Date(),
      gameState: 0,
    };
    this.roundState = {
      roundNumber: 0,
      currentCardCzarId: "",
      blackCard: {
        id: 0,
        cardsToPlay: 0,
        bodyText: "",
      },
      yourHand: [],
    };

    this.onLobbyStateChange = undefined;
    this.onPlayerListChange = undefined;
    this.onRoundStateChange = undefined;
    this.onCommandError = undefined;
    this.onChangeSettings = undefined;
    this.onError = undefined;

    const url = wsBaseUrl;
    console.log(`Connecting to ${url}`);

    /**
     * Mock new WebSocket when testing
     */
    const ws: WebSocket = new WebSocket(url);
    this.wsClient = toWebSocketClient(ws, {
      onDisconnect: () => {
        console.error("Disconnected from the game server");
        this.onError?.("Disconnected from the server");
      },
      onConnect: () => {
        console.log("Connected to the game server");
      },
      onReceive: (msg: string) => {
        this.handleRpcMessage(msg);
      },
      onError: (msg: string) => {
        console.error(`An error ${msg} occurred`);
        // The ping has no body so we don't bother to check it
        this.onError?.(msg);
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

  public emitState() {
    this.onLobbyStateChange?.(structuredClone(this.lobbyState));
    this.onPlayerListChange?.(this.playerList());
    this.onRoundStateChange?.(structuredClone(this.roundState));
  }

  public isOwner(): boolean {
    return this.playerId === this.lobbyState.ownerId;
  }

  private handleOnJoin(msg: RpcOnJoinMsg) {
    const state = msg.state as GameStateInfo;
    this.lobbyState = {
      ownerId: state.gameOwnerId,
      settings: state.settings,
      creationTime: new Date(state.creationTime),
      gameState: state.gameState,
    };

    for (const player of state.players) {
      if (!this.players.find((x) => x.id === player.id)) {
        this.players.push({
          id: player.id,
          name: player.name,
          connected: player.connected,
          points: player.points,
        });
      }
    }

    this.onLobbyStateChange?.(this.lobbyState);
    this.onPlayerListChange?.(this.playerList());
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
    this.lobbyState.ownerId = msg.id;
    this.onLobbyStateChange?.(structuredClone(this.lobbyState));
  }

  private handlePing() {
    this.wsClient?.sendMessage(JSON.stringify(this.encodeMessage(MsgPing, {})));
  }

  private handleRoundInformation(data: RpcRoundInformationMsg) {
    this.roundState = data;
    this.lobbyState.gameState = GameStateWhiteCardsBeingSelected;
    this.onLobbyStateChange?.(structuredClone(this.lobbyState));
    this.onRoundStateChange?.(this.roundState);
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
        this.lobbyState.settings = data.settings;
        this.onLobbyStateChange?.(this.lobbyState);
        return this.onChangeSettings?.(data);
      case MsgOnPlayerLeave:
        console.log("Handling on player leave message");
        return this.handleOnPlayerLeave(rpcMessage.data as RpcOnPlayerLeaveMsg);
      case MsgNewOwner:
        console.log("Handling new owner message");
        return this.handleOnOwnerChange(rpcMessage.data as RpcNewOwnerMsg);
      case MsgPing:
        console.log("Handling ping message");
        return this.handlePing();
      case MsgRoundInformation:
        console.log("Handling round information message");
        return this.handleRoundInformation(
          rpcMessage.data as RpcRoundInformationMsg,
        );
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

    const resp = apiClient.games.leaveDelete();
    this.onError = undefined;
    this.wsClient.disconnect();
    return resp;
  }

  public async startGame() {
    console.log("Starting game: ", this.gameId);
    if (!this.wsClient) {
      throw new Error("Cannot start game as websocket is not connected");
    }

    this.wsClient.sendMessage(
      JSON.stringify(this.encodeMessage(MsgStartGame, {})),
    );
  }

  public playCards(cards: number[]) {
    console.log("Playing cards: ", cards);
    if (!this.wsClient) {
      throw new Error("Cannot play cards as websocket is not connected");
    }

    const msg: RpcPlayCardsMsg = {
      cardIds: cards,
    };

    this.wsClient.sendMessage(
      JSON.stringify(this.encodeMessage(MsgPlayCards, msg)),
    );
  }
}

export const gameState: GameState = new GameState();
