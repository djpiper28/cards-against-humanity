import WebSocket from "ws";

export interface WebSocketClient {
  readonly sendMessage: (msg: string) => void;
  onReceive?: (msg: string) => void;
  onDisconnect?: () => void;
  onConnect?: () => void;
}

export function toWebSocketClient(ws: WebSocket): WebSocketClient {
  const ret: WebSocketClient = {
    sendMessage: ws.send,
  };

  ws.on("close", () => {
    ret?.onDisconnect();
  });
  ws.on("connection", () => {
    ret?.onConnect();
  });
  ws.on("message", (data: string) => {
    ret?.onReceive(data);
  });
  return ret;
}
