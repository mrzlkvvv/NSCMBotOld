build:
	docker-compose build

build-local:
	go build -o ./bin/NSCMBot ./cmd/main.go

run:
	docker-compose up

update:
	go get -u ./...
	go mod tidy
