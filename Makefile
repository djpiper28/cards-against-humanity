all: swagger frontend backend
	echo "Building Done"

swagger:
	go install github.com/swaggo/swag/cmd/swag@latest
	swag init

frontend-install:
	cd ./cahfrontend/ && npm i

frontend-api: swagger
	npx swagger-typescript-api -p ./docs/swagger.json -o ./cahfrontend/src/ -n api.ts

frontend-types:
	go install github.com/gzuidhof/tygo@latest
	tygo generate

frontend-storybook: frontend-install
	cd ./cahfrontend && npm run build-storybook

frontend-main: frontend-install swagger frontend-api frontend-types
	cd ./cahfrontend && npm run build

frontend: frontend-main frontend-storybook
	echo "Building Frontend Done"

backend: swagger
	go build

test-backend: backend
	go test './...'

test-frontend: frontend
	cd ./cahfrontend && npm run test

test: test-backend test-frontend
	echo "Testing Done"

bench: backend
	go test '-bench=./...'

backend-fmt:
	swag fmt
	gofmt -l -w .

frontend-fmt:
	cd ./cahfrontend/ && prettier -w .

fmt: backend-fmt frontend-fmt
	echo "Formatting Done"
