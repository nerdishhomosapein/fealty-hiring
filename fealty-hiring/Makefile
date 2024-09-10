build:
	@go build -o bin/fealty

run: build
	@./bin/fealty
seed: build
	@./bin/fealty -seed

test:
	@go test ./... -v
