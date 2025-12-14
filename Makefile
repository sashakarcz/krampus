.PHONY: help build build-frontend build-backend clean run dev docker-build docker-run docker-stop install test

# Variables
BINARY_NAME=krampus-server
DOCKER_IMAGE=krampus-santa-sync
DOCKER_TAG=latest
CLIENT_DIR=client
SERVER_DIR=server
DIST_DIR=$(CLIENT_DIR)/dist
STATIC_DIR=$(SERVER_DIR)/static

# Default target
.DEFAULT_GOAL := help

help: ## Show this help message
	@echo "Krampus Santa Sync Server - Build Commands"
	@echo "==========================================="
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

install: ## Install all dependencies (Go and Node)
	@echo "Installing Go dependencies..."
	go mod download
	@echo "Installing Node dependencies..."
	cd $(CLIENT_DIR) && npm install
	@echo "All dependencies installed!"

build-frontend: ## Build the React frontend
	@echo "Building frontend..."
	cd $(CLIENT_DIR) && npm run build
	@echo "Frontend built successfully!"

build-backend: ## Build the Go backend only (requires frontend to be built first)
	@echo "Building backend..."
	@if [ ! -d "$(DIST_DIR)" ]; then \
		echo "Error: Frontend not built. Run 'make build-frontend' first."; \
		exit 1; \
	fi
	@echo "Copying frontend to server/static..."
	@rm -rf $(STATIC_DIR)
	@cp -r $(DIST_DIR) $(STATIC_DIR)
	@echo "Building Go binary..."
	go build -ldflags="-s -w" -o $(BINARY_NAME) ./$(SERVER_DIR)
	@echo "Backend built successfully: $(BINARY_NAME)"

build: install build-frontend build-backend ## Build everything (frontend + backend)
	@echo "Complete build finished!"
	@echo "Run with: ./$(BINARY_NAME)"

rebuild: clean build ## Clean and rebuild everything

clean: ## Remove build artifacts
	@echo "Cleaning build artifacts..."
	@rm -rf $(DIST_DIR)
	@rm -rf $(STATIC_DIR)
	@rm -rf $(CLIENT_DIR)/node_modules
	@rm -f $(BINARY_NAME)
	@rm -rf database/krampus.db
	@echo "Clean complete!"

clean-builds: ## Remove only build outputs (keep dependencies)
	@echo "Cleaning build outputs..."
	@rm -rf $(DIST_DIR)
	@rm -rf $(STATIC_DIR)
	@rm -f $(BINARY_NAME)
	@echo "Build outputs cleaned!"

run: ## Run the server (requires build first)
	@if [ ! -f "$(BINARY_NAME)" ]; then \
		echo "Error: Binary not found. Run 'make build' first."; \
		exit 1; \
	fi
	./$(BINARY_NAME)

dev-backend: ## Run backend in development mode (hot reload with air)
	@if ! command -v air &> /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	cd $(SERVER_DIR) && air

dev-frontend: ## Run frontend in development mode
	cd $(CLIENT_DIR) && npm run dev

dev: ## Run both frontend and backend in development mode (requires 2 terminals)
	@echo "To run in dev mode, open 2 terminals:"
	@echo "Terminal 1: make dev-backend"
	@echo "Terminal 2: make dev-frontend"
	@echo "Then access frontend at http://localhost:3000"

test: ## Run all tests
	@echo "Running Go tests..."
	go test -v ./...
	@echo "Running frontend tests..."
	cd $(CLIENT_DIR) && npm test

fmt: ## Format Go code
	go fmt ./...
	gofmt -s -w .

lint: ## Lint Go code
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "Installing golangci-lint..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest; \
	fi
	golangci-lint run

# Docker targets
docker-build: ## Build Docker image
	@echo "Building Docker image (Frontend built inside container)..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "Docker image built: $(DOCKER_IMAGE):$(DOCKER_TAG)"

docker-run: ## Run Docker container
	docker run -d \
		--name krampus \
		-p 8080:8080 \
		-v $(PWD)/database:/app/database \
		-v $(PWD)/.env:/app/.env \
		--restart unless-stopped \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "Container started! Access at http://localhost:8080"
	@echo "View logs: docker logs -f krampus"

docker-run-interactive: ## Run Docker container interactively
	docker run -it --rm \
		-p 8080:8080 \
		-v $(PWD)/database:/app/database \
		-v $(PWD)/.env:/app/.env \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

docker-stop: ## Stop and remove Docker container
	docker stop krampus || true
	docker rm krampus || true

docker-logs: ## View Docker container logs
	docker logs -f krampus

docker-shell: ## Open shell in running container
	docker exec -it krampus /bin/sh

docker-clean: docker-stop ## Remove Docker image and container
	docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) || true
	@echo "Docker artifacts cleaned!"

docker-compose-up: ## Start with docker-compose
	docker-compose up -d
	@echo "Services started! Access at http://localhost:8080"

docker-compose-down: ## Stop docker-compose services
	docker-compose down

docker-compose-logs: ## View docker-compose logs
	docker-compose logs -f

# Database targets
db-init: ## Initialize database (creates tables)
	@echo "Database will be initialized on first run"
	@mkdir -p database

db-clean: ## Delete database (WARNING: destroys all data)
	@echo "Are you sure? This will delete all data. Press Ctrl+C to cancel."
	@sleep 3
	rm -f database/krampus.db
	@echo "Database deleted!"

db-backup: ## Backup database
	@mkdir -p backups
	@cp database/krampus.db backups/krampus-$(shell date +%Y%m%d-%H%M%S).db
	@echo "Database backed up to backups/"

db-shell: ## Open SQLite shell
	sqlite3 database/krampus.db

# Release targets
release: clean build ## Build production release
	@echo "Creating release package..."
	@mkdir -p release
	@cp $(BINARY_NAME) release/
	@cp .env.example release/.env
	@cp README.md release/
	@echo "Release package created in release/"

release-tar: release ## Create release tarball
	@tar -czf krampus-$(shell git describe --tags --always).tar.gz release/
	@echo "Release tarball created: krampus-$(shell git describe --tags --always).tar.gz"

# Utility targets
version: ## Show version information
	@echo "Krampus Santa Sync Server"
	@echo "Go version: $(shell go version)"
	@echo "Node version: $(shell node --version)"
	@echo "npm version: $(shell npm --version)"

check-env: ## Check if .env file exists and is configured
	@if [ ! -f .env ]; then \
		echo "Warning: .env file not found!"; \
		echo "Copy .env.example to .env and configure:"; \
		echo "  cp .env.example .env"; \
		exit 1; \
	fi
	@echo ".env file exists"
	@grep -q "OIDC_PROVIDER_URL=" .env && echo "✓ OIDC_PROVIDER_URL configured" || echo "✗ OIDC_PROVIDER_URL not set"
	@grep -q "OIDC_CLIENT_ID=" .env && echo "✓ OIDC_CLIENT_ID configured" || echo "✗ OIDC_CLIENT_ID not set"
	@grep -q "JWT_SECRET=" .env && echo "✓ JWT_SECRET configured" || echo "✗ JWT_SECRET not set"

setup: ## Initial setup (install deps, create .env, init db)
	@echo "Running initial setup..."
	@if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo "Created .env file from .env.example"; \
		echo "Please edit .env with your configuration"; \
	fi
	@make install
	@make db-init
	@echo ""
	@echo "Setup complete! Next steps:"
	@echo "1. Edit .env with your OIDC credentials"
	@echo "2. Run 'make build' to build the application"
	@echo "3. Run 'make run' to start the server"
