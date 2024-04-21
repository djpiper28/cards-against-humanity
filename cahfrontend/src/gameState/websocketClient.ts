import WebSocket from "isomorphic-ws";

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
  callbacks: WebSocketClientCallbacks,
): WebSocketClient {
  const ret: WebSocketClient = {
    sendMessage: (msg: string) => {
      console.log(`Sending message: ${msg}`);
      ws.send(msg, (err?: Error) => {
        console.error(`Cannot send message ${err}`);
        ws.close();
      });
    },
    ...callbacks,
  };

  ws.onclose = () => {
    ret.onDisconnect();
  };
  ws.onopen = () => {
    ret.onConnect();
  };
  ws.onmessage = function incoming(data) {
    ret.onReceive(data.data.toString());
  };
  ws.onerror = (event) => {
    console.error(event);
    ret.onDisconnect();
    ws.close();
  };
  return ret;
}
