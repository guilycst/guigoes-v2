.PHONY: install build run check rundocker

install:
	go mod download

build: install
	go run ./cmd/blog build

check: build
	go test ./...
	go vet ./...

run: build
	go run ./cmd/blog serve

rundocker:
	docker build -t guigoes-margo:latest .
	docker run -p 8080:80 --rm -it guigoes-margo:latest
