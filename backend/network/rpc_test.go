package network_test

import (
	"encoding/json"
	"testing"

	"github.com/djpiper28/cards-against-humanity/backend/network"
	"github.com/stretchr/testify/assert"
)

type sampleRpcMessage struct {
	Data int
}

type sampleRpcMessageBody struct {
	Data sampleRpcMessage `json:"data"`
}

func (msg sampleRpcMessage) Type() network.RpcMessageType {
	return 123
}

func TestRpcMessageEncode(t *testing.T) {
	t.Parallel()

	data := 123123
	var unencodedMsg network.RpcMessage = sampleRpcMessage{Data: data}
	msg, err := network.EncodeRpcMessage(unencodedMsg)
	assert.Nil(t, err, "Should be able to encode message")

	var rpcMsg sampleRpcMessageBody
	err = json.Unmarshal(msg, &rpcMsg)
	assert.Nil(t, err)

	assert.Equal(t, unencodedMsg, rpcMsg.Data)
}
