#!/bin/sh
cd ../cahfrontend/ || exit 1
export VITE_API_BASE_URL="http://localhost:3255/api"
export VITE_WS_BASE_URL="ws://localhost:3255/ws"
pnpm dev
