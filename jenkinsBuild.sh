#!/bin/sh
GOPATH=/var/lib/jenkins/go
PATH=$PATH:/usr/local/bin
PATH=$PATH:$GOPATH/bin

BASE_URL=cards.djpiper28.co.uk/api
VITE_API_BASE_URL=https://$BASE_URL
VITE_WS_BASE_URL=wss://$BASE_URL/join

make -j
