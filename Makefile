build:
	swag init
	go build

test:
	go test './...'

bench:
	go test '-bench=./...'

fmt:
	gofmt -l -w .
	swag fmt
