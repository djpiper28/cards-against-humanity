#!/bin/sh
export GOPATH=/var/lib/jenkins/go
export PATH=$PATH:/usr/local/bin
export PATH=$PATH:$GOPATH/bin
export PATH=$PATH:/usr/local/go/bin

export BASE_URL=cards.djpiper28.co.uk/api
export VITE_API_BASE_URL=https://$BASE_URL
export VITE_WS_BASE_URL=wss://cards.djpiper28.co.uk/ws

make -j
