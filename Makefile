BINARY := bin/gophertalk

.PHONY: build run vet up down

build:
	go build -o $(BINARY) ./cmd

run:
	go run ./cmd

vet:
	go vet ./...

up:
	docker compose up --build

down:
	docker compose down
