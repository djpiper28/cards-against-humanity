#!/bin/sh
export PATH=/var/lib/jenkins/go/bin:$PATH
export GOPATH=/var/lib/jenkins/go

export VITE_API_BASE_URL=https://cards.djpiper28.co.uk/
export VITE_WS_BASE_URL=wss://cards.djpiper28.co.uk/api/games/join

make frontend -j
