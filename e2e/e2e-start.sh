#!/bin/sh
cd ../cahfrontend/ || exit 1
export VITE_API_BASE_URL="http://localhost:8080"
pnpm dev
