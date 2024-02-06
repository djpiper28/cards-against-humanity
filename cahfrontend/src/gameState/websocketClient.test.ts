// @vitest-environment node
import { v4 } from "uuid";
import { beforeAll, beforeEach, describe, it, expect } from "vitest";
import { toWebSocketClient } from "./websocketClient";
import WebSocket, { WebSocketServer } from "ws";

describe("WebSocketClient tests", () => {
  const port = 1442;
  const onConnectMessage = `${v4()} test message from server`;
  let serverReceivedMessages: string[] = [];
  let wss: WebSocketServer | undefined = undefined;

  beforeAll(() => {
    wss = new WebSocketServer({
      port: port,
    });
    wss.on("connection", (ws) => {
      ws.send(onConnectMessage);
      ws.on("message", (data: string) => {
        serverReceivedMessages.push(data);
      });
    });
  });

  beforeEach(() => {
    wss?.close();
    serverReceivedMessages = [];
  });

  it("Test websocket", () => {
    const ws = new WebSocket(`ws://localhost:${port}`);
    const wsClient = toWebSocketClient(ws);

    let onConnectCalled = false;
    wsClient.onConnect = () => {
      onConnectCalled = true;
    };

    let onDisconnectCalled = false;
    wsClient.onDisconnect = () => {
      onDisconnectCalled = true;
    };

    const clientReceivedMessages: string[] = [];
    wsClient.onReceive = clientReceivedMessages.push;

    const msg = `${v4()} test message from client`;
    wsClient.sendMessage(msg);

    ws.close();
    expect(onConnectCalled).toBe(true);
    expect(onDisconnectCalled).toBe(true);
    expect(serverReceivedMessages).toContain(msg);
    expect(serverReceivedMessages).toHaveLength(1);

    expect(clientReceivedMessages).toContain(onConnectMessage);
    expect(clientReceivedMessages).toHaveLength(1);
  });
});
