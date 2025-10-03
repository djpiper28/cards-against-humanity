// @vitest-environment node
import WebSocket from "isomorphic-ws";
import { v4 } from "uuid";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { wsBaseUrl } from "../apiClient";
import {
  MsgChangeSettings,
  MsgCommandError,
  MsgNewOwner,
  MsgOnBlackCardSkipped,
  MsgOnCardPlayed,
  MsgOnCzarJudgingPhase,
  MsgOnPlayerCreate,
  MsgOnPlayerDisconnect,
  MsgOnPlayerJoin,
  MsgOnPlayerLeave,
  RpcChangeSettingsMsg,
  RpcMessage,
  RpcOnBlackCardSkipped,
  RpcOnCzarJudgingPhaseMsg,
} from "../rpcTypes";
import { gameState } from "./gameState";
import { GameStateCzarJudgingCards, WhiteCard } from "../gameLogicTypes";

describe("Game state tests", () => {
  vi.mock("isomorphic-ws", () => {
    const wsConstructor = vi.fn(function con() {
      return {
        send: vi.fn(),
        onopen: vi.fn(),
        onerror: vi.fn(),
        onclose: vi.fn(),
        onmessage: vi.fn(),
      };
    });

    return { default: wsConstructor };
  });

  beforeEach(() => {
    vi.resetAllMocks();
  });

  it("Should be able to setup a game state", () => {
    expect(gameState).toBeTruthy();
    expect(gameState.validate()).toBeFalsy();

    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");
    expect(gameState.validate()).toBeTruthy();
    expect(WebSocket).toHaveBeenCalledWith(`${wsBaseUrl}`);
  });

  it("Should add a new player to the player list when they connect", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    const msg: RpcMessage = {
      type: MsgOnPlayerJoin,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    expect(gameState.playerList().length).toBe(0);
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList().length).toBe(1);
    expect(gameState.playerList()[0]).toEqual({
      id: msg.data.id,
      name: msg.data.name,
      connected: true,
      points: 0,
      hasPlayed: false,
    });
  });

  it("Should call on player joined method on join", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");
    gameState.onPlayerListChange = vi.fn();

    const msg: RpcMessage = {
      type: MsgOnPlayerJoin,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    expect(gameState.playerList().length).toBe(0);
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList().length).toBe(1);
    expect(gameState.playerList()[0]).toEqual({
      id: msg.data.id,
      name: msg.data.name,
      connected: true,
      points: 0,
      hasPlayed: false,
    });

    expect(gameState.onPlayerListChange).toBeCalledWith(gameState.playerList());
  });

  it("Should not duplicate a player if they are added twice", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    const msg: RpcMessage = {
      type: MsgOnPlayerJoin,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.handleRpcMessage(JSON.stringify(msg));
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList()[0]).toEqual({
      id: msg.data.id,
      name: msg.data.name,
      connected: true,
      points: 0,
      hasPlayed: false,
    });
    expect(gameState.playerList().length).toBe(1);
  });

  it("Should set a player to connected if they are in the players list but disconnected", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    const msg: RpcMessage = {
      type: MsgOnPlayerJoin,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.players = [
      {
        id: msg.data.id,
        name: msg.data.name,
        connected: false,
        points: 0,
        hasPlayed: false,
      },
    ];

    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList().length).toBe(1);
    expect(gameState.playerList()[0]).toEqual({
      id: msg.data.id,
      name: msg.data.name,
      connected: true,
      points: 0,
      hasPlayed: false,
    });
  });

  it("Should add a player when they are created", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    const msg: RpcMessage = {
      type: MsgOnPlayerCreate,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.onPlayerListChange = vi.fn();
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList().length).toBe(1);
    expect(gameState.playerList()[0]).toEqual({
      id: msg.data.id,
      name: msg.data.name,
      connected: false,
      points: 0,
      hasPlayed: false,
    });
    expect(gameState.onPlayerListChange).toBeCalledWith(gameState.playerList());
  });

  it("Should set a player to disconnected when they disconnect", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    // Join as the player
    const joinMsg: RpcMessage = {
      type: MsgOnPlayerCreate,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.handleRpcMessage(JSON.stringify(joinMsg));

    // Disconnect the player
    const msg: RpcMessage = {
      type: MsgOnPlayerDisconnect,
      data: {
        id: joinMsg.data.id,
      },
    };

    gameState.onPlayerListChange = vi.fn();
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList().length).toBe(1);
    expect(gameState.playerList()[0]).toEqual({
      id: msg.data.id,
      name: joinMsg.data.name,
      connected: false,
      points: 0,
      hasPlayed: false,
    });
    expect(gameState.onPlayerListChange).toBeCalledWith(gameState.playerList());
  });

  it("Should encode a RPC message", () => {
    var msg: RpcChangeSettingsMsg = {
      settings: {
        maxPlayers: 4,
        name: "Test Game",
        password: "password",
        public: true,
      },
    };

    var encoded = gameState.encodeMessage(MsgChangeSettings, msg);

    expect(encoded).toEqual({
      type: MsgChangeSettings,
      data: msg,
    });
  });

  it("Should call callback function when a command error occurs", () => {
    const msg: RpcMessage = {
      type: MsgCommandError,
      data: {
        reason: "An error occurred",
      },
    };

    gameState.onCommandError = vi.fn();
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.onCommandError).toBeCalledWith(msg.data);
  });

  it("Should call callback function when settings change", () => {
    const msg: RpcMessage = {
      type: MsgChangeSettings,
      data: {
        settings: {
          maxPlayers: 4,
          name: "Test Game",
          password: "password",
          public: true,
        },
      },
    };

    gameState.onChangeSettings = vi.fn();
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.onChangeSettings).toBeCalledWith(msg.data);
  });

  it("Should send a command and reset the error", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    const msg: RpcMessage = {
      type: MsgChangeSettings,
      data: {
        settings: {
          maxPlayers: 4,
          name: "Test Game",
          password: "password",
          public: true,
        },
      },
    };

    gameState.wsClient = {
      sendMessage: vi.fn(),
    };

    gameState.onCommandError = vi.fn();
    gameState.sendRpcMessage(MsgChangeSettings, msg.data);

    expect(gameState.onCommandError).toBeCalledWith({ reason: "" });
    expect(gameState.wsClient?.sendMessage).toHaveBeenCalledWith(
      JSON.stringify(msg),
    );
  });

  it("Should remove a player from the player list when they disconnect", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    // Join as the player
    const joinMsg: RpcMessage = {
      type: MsgOnPlayerCreate,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.handleRpcMessage(JSON.stringify(joinMsg));

    // Disconnect the player
    const msg: RpcMessage = {
      type: MsgOnPlayerLeave,
      data: {
        id: joinMsg.data.id,
      },
    };

    gameState.onPlayerListChange = vi.fn();
    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList().length).toBe(0);
    expect(gameState.onPlayerListChange).toBeCalledWith(gameState.playerList());
  });

  it("Should emit the state and change the owner id when a new owner is set", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    // Join as the player
    const joinMsg: RpcMessage = {
      type: MsgOnPlayerCreate,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.handleRpcMessage(JSON.stringify(joinMsg));

    const newOwnerMsg: RpcMessage = {
      type: MsgNewOwner,
      data: {
        id: v4(),
      },
    };
    gameState.onLobbyStateChange = vi.fn();

    gameState.handleRpcMessage(JSON.stringify(newOwnerMsg));
    expect(gameState.isOwner()).toBe(false);
    expect(gameState.onLobbyStateChange).toBeCalledTimes(1);

    gameState.playerId = newOwnerMsg.data.id;
    expect(gameState.isOwner()).toBe(true);
  });

  it("Should default to showing that players have not played", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    const joinMsg: RpcMessage = {
      type: MsgOnPlayerCreate,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.handleRpcMessage(JSON.stringify(joinMsg));

    expect(gameState.playerList()[0].hasPlayed).toBe(false);
  });

  it("Should show a player has played after rpc message", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    const joinMsg: RpcMessage = {
      type: MsgOnPlayerCreate,
      data: {
        id: v4(),
        name: "Player 1",
      },
    };

    gameState.handleRpcMessage(JSON.stringify(joinMsg));

    const playedMsg: RpcMessage = {
      type: MsgOnCardPlayed,
      data: {
        playerId: joinMsg.data.id,
      },
    };

    gameState.handleRpcMessage(JSON.stringify(playedMsg));

    expect(gameState.playerList()[0].hasPlayed).toBe(true);
  });

  it("Should handle moving to czar judging phase", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    gameState.onLobbyStateChange = vi.fn();
    gameState.onRoundStateChange = vi.fn();
    gameState.onPlayerListChange = vi.fn();
    gameState.onAllPlaysChanged = vi.fn();

    const allPlays: WhiteCard[][] = [
      [
        {
          id: 1,
          bodyText: "uwu 1",
        },
        {
          id: 2,
          bodyText: "uwu 2",
        },
      ],
      [
        {
          id: 3,
          bodyText: "owo 1",
        },
        {
          id: 4,
          bodyText: "owo 2",
        },
      ],
      [
        {
          id: 5,
          bodyText: "0w0 1",
        },
        {
          id: 6,
          bodyText: "0w0 2",
        },
      ],
    ];
    const newHand: WhiteCard[] = [
      {
        id: 1,
        bodyText: "uwu 1",
      },
      {
        id: 2,
        bodyText: "uwu 2",
      },
      {
        id: 3,
        bodyText: "uwu 3",
      },
      {
        id: 4,
        bodyText: "uwu 4",
      },
      {
        id: 5,
        bodyText: "uwu 5",
      },
      {
        id: 6,
        bodyText: "uwu 6",
      },
      {
        id: 7,
        bodyText: "uwu 7",
      },
    ];

    const data: RpcOnCzarJudgingPhaseMsg = {
      allPlays: allPlays,
      newHand: newHand,
    };

    const czarJudgingMsg: RpcMessage = {
      type: MsgOnCzarJudgingPhase,
      data: data,
    };

    gameState.handleRpcMessage(JSON.stringify(czarJudgingMsg));

    expect(gameState.onLobbyStateChange).toHaveBeenCalled();
    expect(gameState.lobbyState.gameState).toBe(GameStateCzarJudgingCards);

    expect(gameState.onRoundStateChange).toHaveBeenCalled();
    expect(gameState.roundState.yourHand).toEqual(newHand);

    expect(gameState.onAllPlaysChanged).toHaveBeenCalledWith(allPlays);

    for (const player of gameState.playerList()) {
      expect(player.hasPlayed).toBe(false);
    }
  });

  it("Should handle skip black card", () => {
    const gid = v4();
    const pid = v4();
    gameState.setupState(gid, pid, "");

    gameState.onPlayerListChange = vi.fn();
    gameState.onRoundStateChange = vi.fn();

    const newCard: RpcOnBlackCardSkipped = {
      newBlackCard: {
        bodyText: "lorem ipsum",
        cardsToPlay: 1,
        id: 1,
      },
    };
    const msg: RpcMessage = {
      type: MsgOnBlackCardSkipped,
      data: newCard,
    };

    gameState.handleRpcMessage(JSON.stringify(msg));
    expect(gameState.onRoundStateChange).toHaveBeenCalled();
    expect(gameState.onPlayerListChange).toHaveBeenCalled();

    for (const player of gameState.playerList()) {
      expect(player.hasPlayed).toBe(false);
    }
  });
});
