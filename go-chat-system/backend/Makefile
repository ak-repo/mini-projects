# Makefile
.PHONY: help build run test migrate-up migrate-down gen-proto docker-build docker-up docker-down clean

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the application
	go build -o bin/chat-server ./cmd/server

run: ## Run the application
	go run ./cmd/server/main.go

test: ## Run tests
	go test -v -race -coverprofile=coverage.out ./...

migrate-up: ## Run database migrations up
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/chatapp?sslmode=disable" up

migrate-down: ## Run database migrations down
	migrate -path migrations -database "postgresql://postgres:postgres@localhost:5432/chatapp?sslmode=disable" down

migrate-create: ## Create a new migration file (usage: make migrate-create name=migration_name)
	migrate create -ext sql -dir migrations -seq $(name)

gen-proto: ## Generate gRPC code from proto files
	protoc --go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		api/proto/v1/*.proto

docker-build: ## Build Docker image
	docker build -t chat-app:latest .

docker-up: ## Start Docker Compose services
	docker-compose -f deploy/docker-compose.yml up -d

docker-down: ## Stop Docker Compose services
	docker-compose -f deploy/docker-compose.yml down

docker-logs: ## View Docker Compose logs
	docker-compose -f deploy/docker-compose.yml logs -f

clean: ## Clean build artifacts
	rm -rf bin/
	rm -f coverage.out

install-tools: ## Install development tools
	go install github.com/golang/protobuf/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

