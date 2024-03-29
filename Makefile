.PHONY: all
all: frontend backend
	echo "Building Done"

# Swagger defs
swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	cd backend && swag init --requiredByDefault 

# Frontend
frontend-install:
	cd ./cahfrontend/ && pnpm i

frontend-api: swagger 
	npx swagger-typescript-api -p ./backend/docs/swagger.json -o ./cahfrontend/src/ -n api.ts

frontend-tygo:
	go install github.com/gzuidhof/tygo@latest
	cd backend && tygo generate

frontend-types: frontend-tygo frontend-api 
	echo "Generated types"

.PHONY: frontend-storybook
frontend-storybook: frontend-install
	cd ./cahfrontend && pnpm run build-storybook

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
	
GO_TEST_ARGS=-v -benchmem -parallel 16 ./... -covermode=atomic -coverprofile=coverage.out -timeout 60s

.PHONY: test-backend
test-backend: backend
	cd ./backend/ &&	go test './...' ${GO_TEST_ARGS}

.PHONY: test-e2e
test-e2e: frontend test-backend
	cd ./e2e/ && go test './...' ${GO_TEST_ARGS}

.PHONY: test
test: test-backend test-frontend test-e2e frontend-storybook
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
