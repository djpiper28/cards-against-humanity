// @vitest-environment node
import WebSocket, { WebSocketServer } from "../../node_modules/ws/index.js";
import { afterAll, beforeAll, describe, expect, it } from "vitest";
import { toWebSocketClient } from "./websocketClient";
import { createServer } from "http";
import { parse } from "url";
import { v4 } from "uuid";

function delay(ms: number) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

describe("WebSocketClient tests (ws dialup)", () => {
  const port = 1442;
  const onConnectMessage = `${v4()} test message from server`;
  const serverReceivedMessages: string[] = [];
  const wss = new WebSocketServer({
    port: port,
  });

  beforeAll(() => {
    wss.on("connection", (ws: WebSocket) => {
      ws.on("message", (msg: Buffer) => {
        serverReceivedMessages.push(msg.toString());
      });
      ws.on("error", console.error);
      ws.send(onConnectMessage);
    });
  });

  afterAll(() => {
    wss.close();
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

describe("WebSocketClient tests (http upgrade)", () => {
  const port = 1443;
  const onConnectMessage = `${v4()} test message from server`;
  const serverReceivedMessages: string[] = [];
  const wss = new WebSocketServer({
    noServer: true,
  });
  const server = createServer();

  beforeAll(() => {
    wss.on("connection", (ws: WebSocket) => {
      ws.on("message", (msg: Buffer) => {
        console.error("A MESSAGE 1!!1", msg.toString());
        serverReceivedMessages.push(msg.toString());
        console.error(serverReceivedMessages);
      });
      ws.on("error", console.error);
      ws.send(onConnectMessage);
    });
    server.on("upgrade", function upgrade(request, socket, head) {
      if (parse(request.url).pathname === "/foo") {
        wss.handleUpgrade(request, socket, head, function done(ws) {
          wss.emit("connection", ws, request);
        });
      }
    });
    server.listen(port);
  });

  afterAll(() => {
    wss.close();
    server.close();
  });

  it("Test websocket", async () => {
    let onConnectCalled = false;
    let onDisconnectCalled = false;
    const clientReceivedMessages: string[] = [];

    const ws = new WebSocket(`http://localhost:${port}/foo`);
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
