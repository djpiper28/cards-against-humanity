package network_test

import (
	"errors"
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/limits"
	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/stretchr/testify/assert"
)

func TestNewGlobalConnectioManager(t *testing.T) {
	assert.NotNil(t, network.GlobalConnectionManager, "Global connection manager should not be nil")
	assert.NotNil(t, network.GlobalConnectionManager.GameConnectionMap)
}

type MockConnection struct {
	received chan []byte
	sent     chan []byte
	closed   bool
}

func NewMockConnection() *MockConnection {
	return &MockConnection{received: make(chan []byte, limits.ChannelSize),
		sent:   make(chan []byte, limits.ChannelSize),
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
	var mock network.NetworkConnection = NewMockConnection()
	assert.NotNil(t, mock, "Mock init should work")
}
