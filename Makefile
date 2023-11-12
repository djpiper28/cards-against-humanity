
build:
	# Generate Swagger Docs
	swag init
	# Generate the API types for the frontend
	npx swagger-typescript-api -p ./docs/swagger.json -o ./cahfrontend/src/ -n api.ts
	go build
	cd ./cahfrontend && npm i && npm run build

test:
	go test './...'

bench:
	go test '-bench=./...'

fmt:
	gofmt -l -w .
	swag fmt
