import {
  BlackCard,
  GameStateCzarJudgingCards,
  GameStateInfo,
  GameStateInLobby,
  GameStateWhiteCardsBeingSelected,
  Player,
  WhiteCard,
} from "../gameLogicTypes";
import {
  MsgChangeSettings,
  MsgCommandError,
  MsgCzarSelectCard,
  MsgNewOwner,
  MsgOnBlackCardSkipped,
  MsgOnCardPlayed,
  MsgOnCzarJudgingPhase,
  MsgOnGameEnd,
  MsgOnJoin,
  MsgOnPlayerCreate,
  MsgOnPlayerDisconnect,
  MsgOnPlayerJoin,
  MsgOnPlayerLeave,
  MsgOnWhiteCardPlayPhase,
  MsgPing,
  MsgPlayCards,
  MsgRoundInformation,
  MsgSkipBlackCard,
  MsgStartGame,
  MsgKickPlayer,
  RpcChangeSettingsMsg,
  RpcCommandErrorMsg,
  RpcCzarSelectCardMsg,
  RpcMessage,
  RpcMessageBody,
  RpcMessageType,
  RpcNewOwnerMsg,
  RpcOnBlackCardSkipped,
  RpcOnCardPlayedMsg,
  RpcOnCzarJudgingPhaseMsg,
  RpcOnGameEnd,
  RpcOnJoinMsg,
  RpcOnPlayerCreateMsg,
  RpcOnPlayerDisconnectMsg,
  RpcOnPlayerJoinMsg,
  RpcOnPlayerLeaveMsg,
  RpcOnWhiteCardPlayPhase,
  RpcPlayCardsMsg,
  RpcRoundInformationMsg,
  RpcSkipBlackCard,
  RpcKickPlayer,
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
    yourPlays: [],
    totalPlays: 0,
  };

  // Events
  public onLobbyStateChange?: (state?: GameLobbyState) => void;
  public onRoundStateChange?: (state?: RpcRoundInformationMsg) => void;
  public onPlayerListChange?: (players: GamePlayerList) => void;
  public onCommandError?: (error: RpcCommandErrorMsg) => void;
  public onChangeSettings?: (settings: RpcChangeSettingsMsg) => void;
  public onAllPlaysChanged?: (plays: WhiteCard[][]) => void;
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
      yourPlays: [],
      totalPlays: 0,
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

    this.players = [];

    for (const player of state.players) {
      this.players.push({
        id: player.id,
        name: player.name,
        connected: player.connected,
        points: player.points,
        hasPlayed: player.hasPlayed,
      });
    }

    const errorWhiteCard: WhiteCard = {
      id: 0,
      bodyText: "Cannot load this card :(",
    };
    const errorBlackCard: BlackCard = {
      id: 0,
      bodyText: "Cannot load this card :(",
      cardsToPlay: 1,
    };

    const roundState: RpcRoundInformationMsg = {
      yourPlays: state.roundInfo.yourPlays.map(
        (x) =>
          state.roundInfo.yourHand.find((y) => y?.id === x) ?? {
            ...errorWhiteCard,
            id: x,
          },
      ),
      yourHand: state.roundInfo.yourHand.map((x) => x ?? errorWhiteCard),
      roundNumber: state.roundInfo.roundNumber,
      currentCardCzarId: state.roundInfo.czarId,
      blackCard: state.roundInfo.blackCard ?? errorBlackCard,
      totalPlays: state.roundInfo.playersWhoHavePlayed.length,
    };

    this.roundState = roundState;
    state.roundInfo.playersWhoHavePlayed.forEach((pid) => {
      this.handleOnPlayerPlay({ playerId: pid });
    });
    this.onAllPlaysChanged(state.allPlays);
    this.emitState();
  }

  public playerList(): GamePlayerList {
    return [...this.players];
  }

  private handleOnPlayerJoin(msg: RpcOnPlayerJoinMsg) {
    if (!this.players.find((x: Player) => x.id === msg.id)) {
      this.players.push({
        id: msg.id,
        name: msg.id,
        connected: true,
        points: 0,
        hasPlayed: false,
      });
    }

    this.players = this.players.map((x) => {
      if (x.id === msg.id) {
        return {
          ...x,
          connected: true,
        };
      }
      return x;
    });

    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnPlayerCreate(msg: RpcOnPlayerCreateMsg) {
    this.players = this.players.filter((x: Player) => x.id !== msg.id);
    this.players.push({
      id: msg.id,
      name: msg.name,
      connected: false,
      points: 0,
      hasPlayed: false,
    });

    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnPlayerDisconnect(msg: RpcOnPlayerDisconnectMsg) {
    this.players = this.players.map((x) => {
      if (x.id === msg.id) {
        return {
          ...x,
          connected: false,
        };
      }
      return x;
    });

    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnPlayerLeave(msg: RpcOnPlayerLeaveMsg) {
    this.players = this.players.filter((x: Player) => x.id !== msg.id);
    this.onPlayerListChange?.(this.playerList());
  }

  private handleOnPlayerPlay(msg: RpcOnCardPlayedMsg) {
    this.players = this.players.map((x) => {
      if (x.id === msg.playerId) {
        return {
          ...x,
          hasPlayed: true,
        };
      }
      return x;
    });

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
    this.onRoundStateChange?.(structuredClone(this.roundState));

    this.players = this.players.map((player) => {
      return {
        ...player,
        hasPlayed: false,
      };
    });
    this.onPlayerListChange?.(this.players);
  }

  private handleOnCzarJudgingPhase(data: RpcOnCzarJudgingPhaseMsg) {
    this.roundState.yourHand = data.newHand.map((x) => {
      return {
        id: x?.id ?? -1,
        bodyText: x?.bodyText ?? "Cannot load this card :(",
      };
    });
    this.onRoundStateChange?.(structuredClone(this.roundState));

    this.lobbyState.gameState = GameStateCzarJudgingCards;
    this.onLobbyStateChange?.(structuredClone(this.lobbyState));
    this.onAllPlaysChanged?.(
      data.allPlays.map((x) =>
        x.map((card) => {
          return {
            id: card?.id ?? -1,
            bodyText: card?.bodyText ?? "Cannot load this card :(",
          };
        }),
      ),
    );
  }

  public handleOnWhiteCardPlayPhase(msg: RpcOnWhiteCardPlayPhase) {
    this.lobbyState.gameState = GameStateWhiteCardsBeingSelected;

    this.roundState.roundNumber++;
    this.roundState.totalPlays = 0;
    this.roundState.yourPlays = [];
    this.roundState.blackCard = msg.blackCard!;
    this.roundState.yourHand = msg.yourHand as WhiteCard[];
    this.roundState.currentCardCzarId = msg.cardCzarId;

    this.players = this.players.map((player) => {
      if (player.id === msg.winnerId) {
        return {
          ...player,
          points: player.points + 1,
          hasPlayed: false,
        };
      }

      return player;
    });
    this.emitState();
  }

  private handleOnGameEnd(msg: RpcOnGameEnd) {
    if (msg.winnerId) {
      this.players = this.players.map((player) => {
        if (player.id === msg.winnerId) {
          return {
            ...player,
            points: player.points + 1,
            hasPlayed: false,
          };
        }

        return player;
      });
    }

    this.roundState.totalPlays = 0;
    this.roundState.yourPlays = [];
    this.roundState.yourHand = [];
    this.roundState.currentCardCzarId = "not-in-round";

    this.lobbyState.gameState = GameStateInLobby;
    this.emitState();
  }

  private handleOnBlackCardSkipped(msg: RpcOnBlackCardSkipped) {
    this.roundState.blackCard = msg.newBlackCard;
    this.roundState.yourPlays = [];

    this.players = this.players.map((player) => {
      return {
        ...player,
        hasPlayed: false,
      };
    });

    this.onRoundStateChange?.(structuredClone(this.roundState));
    this.onPlayerListChange?.(this.players);
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
      case MsgOnCardPlayed:
        console.log("Handling card played message");
        return this.handleOnPlayerPlay(rpcMessage.data as RpcOnCardPlayedMsg);
      case MsgOnCzarJudgingPhase:
        console.log("Handling czar judging phase message");
        return this.handleOnCzarJudgingPhase(
          rpcMessage.data as RpcOnCzarJudgingPhaseMsg,
        );
      case MsgOnWhiteCardPlayPhase:
        console.log("Handling on white card play phase message");
        return this.handleOnWhiteCardPlayPhase(
          rpcMessage.data as RpcOnWhiteCardPlayPhase,
        );
      case MsgOnGameEnd:
        console.log("Handling on game end message");
        return this.handleOnGameEnd(rpcMessage.data as RpcOnGameEnd);
      case MsgOnBlackCardSkipped:
        console.log("Handling on black card skipped message");
        return this.handleOnBlackCardSkipped(
          rpcMessage.data as RpcOnBlackCardSkipped,
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

    this.wsClient.disconnect();
    this.wsClient = undefined;

    const resp = await apiClient.games.leaveDelete();
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

  public czarSelectCards(cards: number[]) {
    console.log("Czar is selecting cards: ", cards);
    if (!this.wsClient) {
      throw new Error("Cannot play cards as websocket is not connected");
    }

    const msg: RpcCzarSelectCardMsg = {
      cards: cards,
    };

    this.wsClient.sendMessage(
      JSON.stringify(this.encodeMessage(MsgCzarSelectCard, msg)),
    );
  }

  public czarSkipBlackCard() {
    console.log("Czar is skipping black card");
    if (!this.wsClient) {
      throw new Error("Cannot skip cards as websocket is not connected");
    }

    const msg: RpcSkipBlackCard = {};

    this.wsClient.sendMessage(
      JSON.stringify(this.encodeMessage(MsgSkipBlackCard, msg)),
    );
  }

  public czarKickPlayer(playerId: string) {
    console.log("Czar is kicking a player");
    if (!this.wsClient) {
      throw new Error("Cannot skip cards as websocket is not connected");
    }

    const msg: RpcKickPlayer = {
      playerId,
    };

    this.wsClient.sendMessage(
      JSON.stringify(this.encodeMessage(MsgKickPlayer, msg)),
    );
  }
}

export const gameState: GameState = new GameState();
