.SILENT:

run: build
	./bin/NSCMBot

build:
	go build -o ./bin/NSCMBot ./cmd/main.go

run-compose:
	docker-compose up

build-compose:
	docker-compose build

test:
	go test -vet=off ./...

lint:
	golangci-lint run ./... --config=.golangci.yml

update:
	go get -u ./...
	go mod tidy
