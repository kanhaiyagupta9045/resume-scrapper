build:
	@go build -o bin/resume

run: build
	@./bin/resume

test:
	@go test ./...