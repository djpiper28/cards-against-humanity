// @vitest-environment node
import { v4 } from "uuid";
import { toWebSocketClient } from "./websocketClient";
import { beforeAll, describe, expect, it } from "vitest";
import WebSocket, { WebSocketServer } from "../../node_modules/ws/index.js";

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

describe("WebSocketClient tests", () => {
  const port = 1442;
  const onConnectMessage = `${v4()} test message from server`;
  const serverReceivedMessages: string[] = [];

  beforeAll(() => {
    const wss = new WebSocketServer({
      port: port,
    });
    wss.on("connection", (ws: WebSocket) => {
      ws.on("message", (msg: Buffer) => {
        console.error("A MESSAGE 1!!1", msg.toString());
        serverReceivedMessages.push(msg.toString());
        console.error(serverReceivedMessages);
      });
      ws.on("error", console.error);
      ws.send(onConnectMessage);
    });
  });

  it("Test websocket", async () => {
    let onConnectCalled = false;
    let onDisconnectCalled = false;
    const clientReceivedMessages: string[] = [];

    const ws = new WebSocket(`ws://localhost:${port}`);
    const wsClient = toWebSocketClient(ws, {
      onConnect: () => {
        onConnectCalled = true;
        wsClient.sendMessage(msg);
      },
      onDisconnect: () => {
        onDisconnectCalled = true;
      },
      onReceive: (msg: string) => {
        clientReceivedMessages.push(msg);
      },
    });

    const msg = `${v4()} test message from client`;

    await delay(100);
    expect(onConnectCalled).toBe(true);
    expect(clientReceivedMessages).toContain(onConnectMessage);
    expect(serverReceivedMessages).toContain(msg);

    ws.close();
    await delay(100);
    expect(onDisconnectCalled).toBe(true);
  });
});
