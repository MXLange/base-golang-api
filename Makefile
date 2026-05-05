.PHONY: dev

dev:
	go run ./cmd/server/main.go

build-server:
	go build -o server ./cmd/server/main.go