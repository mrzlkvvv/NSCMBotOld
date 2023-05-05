run:
	go run ./cmd/main.go

build:
	go build -o ./bin/NSCMTelegramBot ./cmd/main.go

update:
	go get -u ./...
	go mod tidy
