# cards-against-humanity
Cards Against Humanity website written in Golang.

Live on [changeme](changeme).

[![codecov](https://codecov.io/gh/djpiper28/cards-against-humanity/graph/badge.svg?token=X6YLDCVVLL)](https://codecov.io/gh/djpiper28/cards-against-humanity)
[![Go](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml)

## Dev Stuff

```sh
# Run the tests
# or `gotestsum` if you are cool
go test ./...

# Build and execute
go build
./cards-against-humanity

# or .\cards-against-humanity.exe if you use WinDoze
```

The server will start on `http://localhost:8080`, Prometheus metrics can be found at `/metrics` (server stats), and
game stats on `/game-metrics`. This setup is jank lol, don't question it though.
