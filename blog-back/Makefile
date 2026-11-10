.PHONY: run build tidy

run:
	go run ./cmd/server

build:
	@mkdir -p bin
	go build -o bin/blog-server ./cmd/server

tidy:
	go mod tidy
