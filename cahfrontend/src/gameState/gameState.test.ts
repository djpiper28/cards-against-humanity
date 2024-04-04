// @vitest-environment node
import WebSocket from "isomorphic-ws";
import { v4 } from "uuid";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { wsBaseUrl } from "../apiClient";
import {
  MsgOnPlayerCreate,
  MsgOnPlayerDisconnect,
  MsgOnPlayerJoin,
  RpcMessage,
} from "../rpcTypes";
import { gameState } from "./gameState";

describe("Game state tests", () => {
  vi.mock("isomorphic-ws", () => {
    const wsConstructor = vi.fn(function con() {
      this.send = vi.fn();
      this.onclose = vi.fn();
      this.onopen = vi.fn();
      this.onmessage = vi.fn();
      this.onerror = vi.fn();
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
      },
    ];

    gameState.handleRpcMessage(JSON.stringify(msg));

    expect(gameState.playerList().length).toBe(1);
    expect(gameState.playerList()[0]).toEqual({
      id: msg.data.id,
      name: msg.data.name,
      connected: true,
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
    });
    expect(gameState.onPlayerListChange).toBeCalledWith(gameState.playerList());
  });
});
