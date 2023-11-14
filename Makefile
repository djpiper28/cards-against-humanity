swagger:
	# Generate Swagger Docs
	swag init
	# Generate the API types for the frontend
	npx swagger-typescript-api -p ./docs/swagger.json -o ./cahfrontend/src/ -n api.ts

frontend: swagger
	cd ./cahfrontend && npm i && npm run build

backend: swagger
	go build

all: swagger frontend backend
	echo "Done"

test: all
	go test './...'

bench: all
	go test '-bench=./...'

fmt:
	gofmt -l -w .
	swag fmt
	cd ./cahfrontend/ && prettier -w .
