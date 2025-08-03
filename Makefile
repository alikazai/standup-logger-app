include .env
export $(shell sed 's/=.*//' .env)

# Build the Standup Logger Application
build:
	@echo "Building application..."
	@go build -o bin/app/standup-logger ./cmd/app

run: build
	@echo "Running application..."
	@./bin/app/standup-logger

# Run migrations
migrate:
	migrate -path db/migrations -database "$(DATABASE_URL)" up

# Rollback latest migration
rollback:
	migrate -path db/migrations -database "$(DATABASE_URL)" down 1

# Export schema to schema.sql
export-schema:
	pg_dump --schema-only --no-owner --file=db/schema.sql "$(DATABASE_URL)"

# Run SQL seed files (requires go run setup)
seed:
	go run cmd/seed/main.go

# Full workflow: migrate, export schema, seed (optional)
# setup-db: migrate export-schema seed
setup-db: migrate seed
	@echo "Running application..."

test:
	@echo "Testing Standup Logger Application"
	@gotestsum --format testname --hide-summary=skipped
	# @go test -v ./... | grep -v 'SKIP\|TODO\|RUN'

test-nocache:
	@echo "Testing Standup Logger Application With No Cache"
	@go test -v ./... -count=1 | grep -v 'SKIP\|TODO\|RUN'

clean:
	@echo "Cleaning..."
	rm -rf ./bin

# Run Docker Compose (up the services defined in docker-compose.yml)
docker:
	@echo "Running Docker Compose..."
	@docker compose up --build -d

# Build Docker image without running containers
docker-build:
	@echo "Building Docker image..."
	@docker compose build

# Stop Docker Compose services
docker-stop:
	@echo "Stopping Docker Compose services..."
	@docker compose down

# Clean up Docker images (optional)
docker-clean:
	@echo "Removing Docker images..."
	@docker compose down --rmi all --remove-orphans

# Connect to Docker container 
docker-connect:
	@echo "Connecting to Docker standup-logger-app container..."
	@docker exec -it standup-logger-app sh

# Create Docker network
docker-network:
	@echo "Creating Docker standup-logger-net network..."
	@docker network create standup-logger-net
