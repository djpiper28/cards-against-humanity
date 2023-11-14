all: swagger frontend backend
	# Done

swagger:
	swag init

frontendapi: swagger
	npx swagger-typescript-api -p ./docs/swagger.json -o ./cahfrontend/src/ -n api.ts

frontendtypes:
	tygo generate

frontend: swagger frontendapi frontendtypes
	cd ./cahfrontend && npm i && npm run build

backend: swagger
	go build

test: all
	go test './...'

bench: all
	go test '-bench=./...'

fmt:
	gofmt -l -w .
	swag fmt
	cd ./cahfrontend/ && prettier -w .
