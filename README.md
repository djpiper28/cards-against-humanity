# cards-against-humanity

Cards Against Humanity website written in Golang.

(WIP) Live on [cards.djpiper28.co.uk](https://cards.djpiper28.co.uk).

[![codecov](https://codecov.io/gh/djpiper28/cards-against-humanity/graph/badge.svg?token=X6YLDCVVLL)](https://codecov.io/gh/djpiper28/cards-against-humanity)
[![Go](https://github.com/djpiper28/cards-against-humanity/actions/workflows/coverage.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/coverage.yml)

## Dev Stuff

### Building

#### Deps

- pnpm
- Go
- GNU Make

To build/test all parts of the system use the below commands.

```sh
# Build and execute
make -j

./backend
# or .\backend.exe if you use WinDoze

# Format the code
make fmt -j

# Run the tests
make test -j

# Run the benchmarks
make bench -j
```

### Backend

The backend is in Go and uses Gin and Gorilla.

The server will start on `http://localhost:8080`, Prometheus metrics can be found at `/metrics` (server stats), and
game stats on `/game-metrics`. This setup is jank lol, don't question it though. The swagger docs can be found on
`http://localhost:8080/swagger/index.html`.

#### Debugging

To debug you can use `gdb` (of lldb if you are on mac), the debug-xxx targets in the Makefile will build end exec the debugger, mac and windows people need to use their debugger I guess.

### Frontend

The frontend is in TS and uses SolidJS, Vite, Vitest, and Storybook. You should use the
Makefile for building `make -j`, and `pnpm dev` for a dev server.

The server will start on `http://localhost:3000`, currently the backend is set to
`http://localhost:8080`
