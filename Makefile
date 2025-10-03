.PHONY: all
all: frontend backend
	echo "Building Done"

# Swagger defs
.PHONY: swagger
swagger:
	cd backend && go generate ./...

# Frontend
.PHONY: frontend-install
frontend-install:
	cd ./cahfrontend/ && pnpm i

.PHONY: frontend-api
frontend-api: swagger 
	npx swagger-typescript-api@13.2.13 generate -p ./backend/docs/swagger.json -o ./cahfrontend/src/ -n api.ts

.PHONY: frontend-tygo
frontend-tygo:
	cd backend && go run github.com/gzuidhof/tygo generate

.PHONY: frontend-types
frontend-types: frontend-tygo frontend-api 
	echo "Generated types"

.PHONY: frontend
frontend: frontend-install frontend-types
	cd ./cahfrontend && pnpm run build

# Backend
.PHONY: backend
backend: swagger
	cd ./backend/ && go build

# Tests
.PHONY: test-frontend
test-frontend: frontend-types
	cd ./cahfrontend && pnpm run test
	
GO_TEST_ARGS=-v -benchmem ./... -race -cover -coverpkg ./... -covermode=atomic -coverprofile=coverage.out

.PHONY: test-backend
test-backend: backend
	cd ./backend/ &&	go test './...' ${GO_TEST_ARGS}

.PHONY: start-docker-compose
start-docker-compose:
	docker compose up --build --detach

# Everything is tested within docker, so this can be started imediately
.PHONY: test-e2e
test-e2e: start-docker-compose
	cd ./e2e/ && go test './...' ${GO_TEST_ARGS}

.PHONY: test
test: test-backend test-frontend test-e2e
	echo "Testing Done"

.PHONY: bench
bench: backend
	cd ./backend/ && go test '-bench=./...'

# Formatters
.PHONY: e2e-fmt
e2e-fmt:
	cd ./e2e/ && gofmt -l -w .

.PHONY: backend-fmt
backend-fmt:
	cd ./backend/ && swag fmt && gofmt -l -w .

.PHONY: frontend-fmt
frontend-fmt:
	cd ./cahfrontend/ && prettier -w .

.PHONY: fmt
fmt: backend-fmt frontend-fmt e2e-fmt
	echo "Formatting Done"

# Debug scripts
.PHONY: debug-e2e-tests
debug-e2e-tests:
	cd ./e2e/ && go test -c && gdb ./e2e.test

.PHONY: debug-backend-tests
debug-backend-tests:
	cd ./backend/ && go test -c && gdb ./backend.test

.PHONY: debug-backend
debug-backend:
	cd ./backend/ && go build && gdb ./backend
