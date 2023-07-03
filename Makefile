build:
	docker-compose build

run:
	docker-compose up

update:
	go get -u ./...
	go mod tidy
