export GOWORK := off
export WRANGLER_SEND_METRICS := false

.PHONY: install build run check dry-run verify rundocker

install:
	go mod download

build: install
	go run ./cmd/blog build

check: build
	test -z "$$(gofmt -l cmd internal)"
	go test -count=1 ./...
	go vet ./...

run: build
	go run ./cmd/blog serve

rundocker:
	docker build -t guigoes-margo:latest .
	docker run -p 8080:80 --rm -it guigoes-margo:latest

# Package the validated site without publishing it.
dry-run: build
	npx --no-install wrangler deploy --dry-run --outdir .build/wrangler-output

# Inspect the existing publication on the configured public origin.
verify:
	go run ./cmd/blog verify -timeout 3m
