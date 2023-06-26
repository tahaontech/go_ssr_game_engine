build:
	@go build -o bin/output cmd/main.go

run: build
	@bin/output