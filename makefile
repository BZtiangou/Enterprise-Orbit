.PHONY: build run dev migrate test clean

build:
	@echo "Building Enterprise-Orbit..."
	go build -o bin/enterprise-orbit cmd/server/main.go

run: build
	@echo "Starting Enterprise-Orbit..."
	./bin/enterprise-orbit

dev:
	@echo "Starting in development mode..."
	air

migrate:
	@echo "Running database migrations..."
	go run cmd/server/main.go migrate

test:
	@echo "Running tests..."
	go test ./...

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	go clean

docker-up:
	docker compose up -d

docker-down:
	docker compose down