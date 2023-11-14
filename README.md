# cards-against-humanity
Cards Against Humanity website written in Golang.

Live on [changeme](changeme).

[![codecov](https://codecov.io/gh/djpiper28/cards-against-humanity/graph/badge.svg?token=X6YLDCVVLL)](https://codecov.io/gh/djpiper28/cards-against-humanity)
[![Go](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml)

## Dev Stuff

### Backend

The backend is in Go and uses Gin and Gorilla.

```sh
# Build and execute
make -j

./cards-against-humanity
# or .\cards-against-humanity.exe if you use WinDoze

# Format the code 
make fmt

# Run the tests
make test

# Run the benchmarks
make bench
```

> Setting Up Code Generators

```sh
go install github.com/swaggo/swag/cmd/swag@latest
go install github.com/gzuidhof/tygo@latest
```

The server will start on `http://localhost:8080`, Prometheus metrics can be found at `/metrics` (server stats), and
game stats on `/game-metrics`. This setup is jank lol, don't question it though.

### Frontend

The frontend is in TS and uses SolidJS and Vite is.

```sh
cd ./cahfrontend/
npm i 

# Building
npm run build 
# Output in /dist

# Storybook
npm run storybook
npm run build-storybook
```
