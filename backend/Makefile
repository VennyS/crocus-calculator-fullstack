.PHONY: run migrate test test-integration docker-check clean

GO=go
TEST_DIR=./...
TEST_FLAGS=-v -count=1
ENV=test.env

run:
	go run ./cmd/server/main.go

migrate:
	go run ./cmd/server/main.go --migrate

test:
	@echo "🔬 Running tests with Docker..."
	@export TESTCONTAINERS_PROVIDER=docker && $(GO) test $(TEST_DIR) $(TEST_FLAGS)

test-integration:
	@echo "🧪 Running integration tests only..."
	@export TESTCONTAINERS_PROVIDER=docker && $(GO) test ./internal/storage/... $(TEST_FLAGS)

docker-check:
	@echo "🐳 Checking Docker daemon..."
	@docker info > /dev/null && echo "✅ Docker is running" || echo "❌ Docker is not running"

clean:
	@echo "🧹 Cleaning test cache..."
	@$(GO) clean -testcache

