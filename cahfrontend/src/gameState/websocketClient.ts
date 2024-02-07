import WebSocket from "ws";

export interface WebSocketClientCallbacks {
  readonly onReceive: (msg: string) => void;
  readonly onDisconnect: () => void;
  readonly onConnect: () => void;
}

interface WebSocketSend {
  readonly sendMessage: (msg: string) => void;
}

export type WebSocketClient = WebSocketSend & WebSocketClientCallbacks;

export function toWebSocketClient(
  ws: WebSocket,
  callbacks: WebSocketClientCallbacks
): WebSocketClient {
  const ret: WebSocketClient = {
    sendMessage: (msg: string) => {
      ws.send(msg, console.error);
    },
    ...callbacks,
  };

  ws.on("close", () => {
    ret.onDisconnect();
  });
  ws.on("open", () => {
    ret.onConnect();
  });
  ws.on("message", (buf: Buffer) => {
    ret.onReceive(buf.toString());
  });
  ws.on("error", (err: Error) => {
    console.error(err);
    ret.onDisconnect();
    ws.close();
  });
  return ret;
}
