# Database configuration
DB_HOST := localhost
DB_PORT := 5432
DB_USER := loco_admin
DB_PASSWORD := toortoor
DB_NAME := loco
DB_URL := postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Migration commands
.PHONY: migrate-up migrate-down migrate-status migrate-force migrate-create db-reset db-setup

migrate-up:
	migrate -path ./migrations -database "$(DB_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DB_URL)" down 1

migrate-down-all:
	migrate -path ./migrations -database "$(DB_URL)" down

migrate-status:
	migrate -path ./migrations -database "$(DB_URL)" version

migrate-force:
	@read -p "Enter version to force: " version; \
	migrate -path ./migrations -database "$(DB_URL)" force $$version

migrate-create:
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir ./migrations -seq $$name

# Database setup
db-reset: db-drop db-create migrate-up

db-create:
	createdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) || true

db-drop:
	dropdb -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) $(DB_NAME) --if-exists

db-setup:
	@echo "Setting up database..."
	make db-create
	make migrate-up
	@echo "Database setup complete!"

# Verification
db-verify:
	@echo "Verifying database setup..."
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "\dt"
	@echo "\nLanguages:"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "SELECT language_id, name FROM languages;"
	@echo "\nTags count:"
	@psql -h $(DB_HOST) -p $(DB_PORT) -U $(DB_USER) -d $(DB_NAME) -c "SELECT COUNT(*) FROM tags;"

# Testing
test:
	go test ./... -v

test-coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out

# Build
build:
	go build -o bin/loco ./cmd/server

# Run
run:
	go run ./cmd/server

# Clean
clean:
	rm -rf bin/
	rm -f coverage.out

# Install tools
install-tools:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Linting
lint:
	golangci-lint run
