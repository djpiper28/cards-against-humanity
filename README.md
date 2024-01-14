# cards-against-humanity
Cards Against Humanity website written in Golang.

Live on [changeme](changeme).

[![codecov](https://codecov.io/gh/djpiper28/cards-against-humanity/graph/badge.svg?token=X6YLDCVVLL)](https://codecov.io/gh/djpiper28/cards-against-humanity)
[![Go](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml)

## Dev Stuff

### Deps
 - pnpm
 - Go
 - GNU Make

### Backend

The backend is in Go and uses Gin and Gorilla.

```sh
# Build and execute
make -j

./backend
# or .\backend.exe if you use WinDoze

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
game stats on `/game-metrics`. This setup is jank lol, don't question it though. The swagger docs can be found on 
`http://localhost:8080/swagger/index.html`.

### Frontend

The frontend is in TS and uses SolidJS, Vite, Vitest, and Storybook. You should use the 
Makefile for building `make -j`, and `pnpm dev` for a dev server.
