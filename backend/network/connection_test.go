package network_test

import (
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

const wsIntegrationTestPort = 6655

var wsIntrationTestUrl = fmt.Sprintf("ws://0.0.0.0:%d", wsIntegrationTestPort)
var gid = uuid.New()
var pid = uuid.New()

type MockGlobalConnectionManager struct {
	newConnectionCalled          bool
	calledGameId, calledPlayerId uuid.UUID
	calledConn                   *websocket.Conn
	network.IntegratedConnectionManager
}

func (gcm *MockGlobalConnectionManager) NewConnection(conn *websocket.Conn, gameId, playerId uuid.UUID) *network.WsConnection {
	defer conn.Close()
	gcm.calledGameId = gameId
	gcm.calledPlayerId = playerId
	gcm.calledConn = conn
	return nil
}

type WsUpgradeSuite struct {
	gcm *MockGlobalConnectionManager
	suite.Suite
}

func TestWsUpgradeSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(WsUpgradeSuite))
}

func (s *WsUpgradeSuite) TestConnectingShouldCallNewConnection() {
	t := s.T()

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = time.Millisecond * 100

	log.Print("Dialing server")
	conn, _, err := dialer.Dial(wsIntrationTestUrl, nil)
	assert.Nil(t, err, "Should connect to server without error")
	defer conn.Close()
	assert.NotNil(t, conn, "Should not have a nil connection")

	log.Print("Waiting for server to do its stuff")
	time.Sleep(time.Millisecond * 100)
	assert.Equal(t, pid, s.gcm.calledPlayerId)
	assert.Equal(t, gid, s.gcm.calledGameId)
	assert.NotNil(t, s.gcm.calledConn, "Should have a valid connection on the server")
}

func (s *WsUpgradeSuite) SetupSuite() {
	s.gcm = &MockGlobalConnectionManager{}

	log.Print("Starting ws server")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		network.WsUpgrade(w, r, gid, pid, s.gcm)
	})
	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", wsIntegrationTestPort), nil)

	time.Sleep(time.Millisecond * 100)
	log.Print("Ws server started")
}
