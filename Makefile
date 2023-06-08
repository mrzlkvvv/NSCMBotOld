MAIN_GO_PATH = ./cmd/main.go
BOT_BIN_PATH = ./bin/NSCMTelegramBot

run:
	go run $(MAIN_GO_PATH)

build:
	go build -o $(BOT_BIN_PATH) $(MAIN_GO_PATH)

fmt:
	gofmt -s -w ./

update:
	go get -u ./...
	go mod tidy

precommit: fmt build
