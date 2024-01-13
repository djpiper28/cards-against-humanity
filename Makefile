all: swagger frontend backend
	echo "Building Done"

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	cd backend && swag init

frontend-install:
	cd ./cahfrontend/ && pnpm i

frontend-api: swagger
	npx swagger-typescript-api -p ./backend/docs/swagger.json -o ./cahfrontend/src/ -n api.ts

frontend-types:
	go install github.com/gzuidhof/tygo@latest
	cd backend && tygo generate

frontend-storybook: frontend-install
	cd ./cahfrontend && pnpm run build-storybook

frontend-main: frontend-install swagger frontend-api frontend-types
	cd ./cahfrontend && pnpm run build

frontend: frontend-main frontend-storybook
	echo "Building Frontend Done"

test-frontend: frontend
	cd ./cahfrontend && pnpm run test

backend: swagger
	cd ./backend/ && go build

test-backend: backend
	cd ./backend/ && go test './...'

test-e2e: frontend backend
	cd ./e2e/ && go test './...'

test: test-backend test-frontend test-e2e
	echo "Testing Done"

bench: backend
	cd ./backend/ && go test '-bench=./...'

e2e-fmt:
	cd ./e2e/ && gofmt -l -w .

backend-fmt:
	cd ./backend/ && swag fmt && gofmt -l -w .

frontend-fmt:
	cd ./cahfrontend/ && prettier -w .

fmt: backend-fmt frontend-fmt e2e-fmt
	echo "Formatting Done"
