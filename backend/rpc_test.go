package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sampleRpcMessage struct {
	Data int
}

type sampleRpcMessageBody struct {
	Data sampleRpcMessage `json:"data"`
}

func (msg sampleRpcMessage) Type() RpcMessageType {
	return 123
}

func TestRpcMessageEncode(t *testing.T) {
	t.Parallel()

	data := 123123
	var unencodedMsg RpcMessage = sampleRpcMessage{Data: data}
	msg, err := EncodeRpcMessage(unencodedMsg)
	assert.Nil(t, err, "Should be able to encode message")

	var rpcMsg sampleRpcMessageBody
	err = json.Unmarshal(msg, &rpcMsg)
	assert.Nil(t, err)

	assert.Equal(t, unencodedMsg, rpcMsg.Data)
}
