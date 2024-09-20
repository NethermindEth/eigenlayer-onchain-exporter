DEFAULT_GOAL := help 
.PHONY: env

run:
	go run main.go run

build:
	go build -o ./bin/eoe

build-docker:
	docker build -t eoe .

pre-commit:
	pre-commit run --all-files
