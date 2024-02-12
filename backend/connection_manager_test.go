package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	conn := NewGameConnection()
	assert.NotNil(t, conn, "Game connection should not be nil")
	conn.CloseAll()
}

func TestNewGlobalConnectioManager(t *testing.T) {
	assert.NotNil(t, globalConnectionManager, "Global connection manager should not be nil")
	assert.NotNil(t, globalConnectionManager.GameConnectionMap)
}

type MockConnection struct {
	received chan []byte
	sent     chan []byte
	closed   bool
}

func NewMockConnection() *MockConnection {
	return &MockConnection{received: make(chan []byte),
		sent:   make(chan []byte),
		closed: false}
}

func (c *MockConnection) Close() {
	c.closed = true
}

func (c *MockConnection) Send(message []byte) error {
	c.sent <- message
	return nil
}

func (c *MockConnection) Receive() ([]byte, error) {
	select {
	case msg := <-c.received:
		return msg, nil
	default:
		return nil, errors.New("queue is kil")
	}
}

func TestSanityCheckMock(t *testing.T) {
	var mock NetworkConnection = NewMockConnection()
	assert.NotNil(t, mock, "Mock init should work")
}
