import WebSocket from "isomorphic-ws";

export interface WebSocketClientCallbacks {
  readonly onReceive: (msg: string) => void;
  readonly onDisconnect: () => void;
  readonly onError: (msg: string) => void;
  readonly onConnect: () => void;
}

interface WebSocketSend {
  readonly sendMessage: (msg: string) => void;
  readonly disconnect: () => void;
}

export type WebSocketClient = WebSocketSend & WebSocketClientCallbacks;

export function toWebSocketClient(
  ws: WebSocket,
  callbacks: WebSocketClientCallbacks,
): WebSocketClient {
  const ret: WebSocketClient = {
    sendMessage: (msg: string) => {
      ws.send(msg, (err?: Error) => {
        console.error(`Cannot send message ${err}`);
        ws.close();
      });
    },
    disconnect: () => {
      ws.close();
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
    ret.onError(event.message);
    ret.onDisconnect();
    ws.close();
  };
  return ret;
}
