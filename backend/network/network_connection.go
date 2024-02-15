package network

import (
	"github.com/gorilla/websocket"
)

// An interface to abstract a connection over a network to allow for networking to
// be tested with mocked connections.
type NetworkConnection interface {
	// A blocking call that will write to the network connection
	Send(message []byte) error
	// A blocking call that will read from the network connection
	Receive() ([]byte, error)
	Close()
}

type WebsocketConnection struct {
	Conn *websocket.Conn
}

func (ws *WebsocketConnection) Send(message []byte) error {
	return ws.Conn.WriteMessage(websocket.TextMessage, []byte(message))
}

func (ws *WebsocketConnection) Receive() ([]byte, error) {
	/*msgType*/ _, msg, err := ws.Conn.ReadMessage()
	// if msgType == websocket.TextMessage {
	return msg, err
}

func (ws *WebsocketConnection) Close() {
	ws.Conn.Close()
}
