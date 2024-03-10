#!/bin/sh
export GOPATH=/var/lib/jenkins/go
export PATH=$PATH:$GOPATH/bin

export VITE_API_BASE_URL=https://cards.djpiper28.co.uk/
export VITE_WS_BASE_URL=wss://cards.djpiper28.co.uk/api/games/join

make frontend -j
