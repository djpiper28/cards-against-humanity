# cards-against-humanity

Cards Against Humanity website written in Golang.

(WIP) Live on [cards.djpiper28.co.uk](https://cards.djpiper28.co.uk).

[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=djpiper28_cards-against-humanity&metric=alert_status)](https://sonarcloud.io/summary/new_code?id=djpiper28_cards-against-humanity)
[![Coverage](https://sonarcloud.io/api/project_badges/measure?project=djpiper28_cards-against-humanity&metric=coverage)](https://sonarcloud.io/summary/new_code?id=djpiper28_cards-against-humanity)
[![Tests](https://github.com/djpiper28/cards-against-humanity/actions/workflows/tests.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/tests.yml)
[![Coverage](https://github.com/djpiper28/cards-against-humanity/actions/workflows/coverage.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/coverage.yml)
[![e2e](https://github.com/djpiper28/cards-against-humanity/actions/workflows/e2e.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/e2e.yml)

## Dev Stuff

You can run the software via Docker Compose, or build it yourself. You can then go to `http://localhost:8000`

```sh
docker-compose up
```

### Building

#### Deps

- pnpm
- Go
- GNU Make
- Docker
- Docker Compose

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

To debug you can use `gdb` (of lldb if you are on mac), mac and windows people need to use their debugger I guess.

### Frontend

The frontend is in TS and uses SolidJS, Vite, Vitest, and Storybook. You should use the
Makefile for building `make -j`, and `pnpm dev` for a dev server.

The server will start on `http://localhost:3000`, currently the backend is set to
`http://localhost:8080`

### Local Development

Due to CORS errors, you need to setup a proxy, see [devProxy](./devProxy/README.md) for a local development proxy server.

Following this you need to run the proxy, backend, and the frontend dev server, then go to the proxies: [`http://localhost:3255`](http://localhost:3255).
