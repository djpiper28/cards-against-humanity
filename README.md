# cards-against-humanity
Cards Against Humanity website written in Golang.

Live on [changeme](changeme).

[![codecov](https://codecov.io/gh/djpiper28/cards-against-humanity/graph/badge.svg?token=X6YLDCVVLL)](https://codecov.io/gh/djpiper28/cards-against-humanity)
[![Go](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml/badge.svg)](https://github.com/djpiper28/cards-against-humanity/actions/workflows/go.yml)

## Dev Stuff

### Backend

The backend is in Go and uses Gin and Gorilla.

```sh
# Run the tests
# or `gotestsum` if you are cool
go test ./...

# Build and execute
make
./cards-against-humanity

# or .\cards-against-humanity.exe if you use WinDoze
```

> Setting Up Swag To Generate Swagger Docs

```sh
go install github.com/swaggo/swag/cmd/swag@latest
swag i # this should now work
```

The server will start on `http://localhost:8080`, Prometheus metrics can be found at `/metrics` (server stats), and
game stats on `/game-metrics`. This setup is jank lol, don't question it though.

### Frontend

The frontend is in TS and uses SolidJS and whatever the f&ck a Vite is. (I'm not a frontend dev, I am a free man)
