# Makefile

.PHONY: proto build run clean test

# Generate proto files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		proto/user_service.proto

# Build the application
build:
	go build -o bin/server main.go

# Run the application
run:
	go run main.go

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate proto and build
all: proto build

# Clean generated files
clean:
	rm -rf bin/
	rm -f users.db

# Run tests
test:
	go test -v ./...

# Install protoc plugins (run once)
install-proto:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Format code
fmt:
	go fmt ./...

# Lint code
lint:
	golangci-lint run