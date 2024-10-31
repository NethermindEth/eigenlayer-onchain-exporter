DEFAULT_GOAL := help 
.PHONY: env

run:
	go run cmd/eoe/main.go run

build:
	go build -o ./bin/eoe cmd/eoe/main.go

build-docker:
	docker build -t eoe .

pre-commit:
	pre-commit run --all-files
