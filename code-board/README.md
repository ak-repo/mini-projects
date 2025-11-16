# code-board

CodeBoard — A developer collaboration platform that combines Git hosting, Kanban boards, and real-time team workflow in one place. Built with Go (Gin) and React

# code-board

## CodeBoard — A developer collaboration platform that combines Git hosting, Kanban boards, and real-time team workflow in one place. Built with Go (Gin) and React


## github.com/ak-repo/code-board
---

# Monorepo Structure

codeboard/
│
├── frontend/ # React frontend
│ ├── codeboard-ui/ # main frontend app
│ ├── package.json
│ └── ...
│
├── gateway/ # 🌐 API Gateway (Gin)
│ ├── cmd/
│ │ └── main.go
│ ├── internal/
│ │ ├── routes/
│ │ │ └── router.go
│ │ ├── middleware/
│ │ │ ├── auth.go # JWT validation, RBAC
│ │ │ ├── cors.go
│ │ │ └── rate_limit.go
│ │ └── proxy/
│ │ └── reverse_proxy.go # forwards to microservices
│ ├── go.mod
│ └── go.sum
│
├── services/
│ ├── trello-service/
│ │ ├── cmd/main.go
│ │ ├── internal/
│ │ │ ├── model/
│ │ │ │ ├── board.go
│ │ │ │ ├── card.go
│ │ │ │ └── comment.go
│ │ │ ├── handler/
│ │ │ │ └── board_handler.go
│ │ │ ├── repository/
│ │ │ │ └── board_repo.go
│ │ │ ├── service/
│ │ │ │ └── board_service.go
│ │ │ ├── events/
│ │ │ │ └── producer.go # Publishes to Kafka
│ │ │ └── db/
│ │ │ └── postgres.go
│ │ ├── go.mod
│ │ └── go.sum
│ │
│ ├── git-service/
│ │ ├── cmd/main.go
│ │ ├── internal/{model,handler,repo,service,events,db}
│ │ ├── integrations/
│ │ │ ├── github.go
│ │ │ └── gitlab.go
│ │ ├── go.mod
│ │ └── go.sum
│ │
│ ├── user-auth-service/
│ │ ├── cmd/main.go
│ │ ├── internal/{model,handler,repo,service,db}
│ │ ├── pkg/jwt/
│ │ │ └── jwt.go
│ │ ├── pkg/rbac/
│ │ │ └── roles.go
│ │ ├── go.mod
│ │ └── go.sum
│ │
│ ├── notification-service/
│ │ ├── cmd/main.go
│ │ ├── internal/
│ │ │ ├── ws/ # WebSocket server
│ │ │ ├── email/
│ │ │ │ └── email_client.go
│ │ │ └── consumer/
│ │ │ └── kafka_consumer.go
│ │ └── go.mod
│ │
│ ├── audit-trail-service/
│ │ ├── cmd/main.go
│ │ ├── internal/{model,handler,repo,service,consumer}
│ │ └── go.mod
│ │
│ ├── search-service/
│ │ ├── cmd/main.go
│ │ ├── internal/
│ │ │ ├── elastic/
│ │ │ │ └── client.go
│ │ │ ├── indexer/
│ │ │ │ └── card_indexer.go
│ │ │ └── query/
│ │ │ └── search.go
│ │ └── go.mod
│ │
│ ├── cache-service/
│ │ ├── cmd/main.go
│ │ ├── internal/redis/
│ │ │ └── client.go
│ │ ├── internal/service/
│ │ │ └── cache_service.go
│ │ └── go.mod
│
├── shared/ # 🧩 Shared libraries for all services
│ ├── pkg/
│ │ ├── config/
│ │ │ └── config.go # Load env vars, config structs
│ │ ├── db/
│ │ │ └── postgres.go
│ │ ├── kafka/
│ │ │ ├── producer.go
│ │ │ └── consumer.go
│ │ ├── logger/
│ │ │ └── logger.go
│ │ ├── middleware/
│ │ │ ├── auth.go
│ │ │ └── cors.go
│ │ ├── response/
│ │ │ └── api_response.go
│ │ ├── utils/
│ │ │ ├── hash.go
│ │ │ └── env.go
│ │ └── dto/
│ │ ├── user_dto.go
│ │ ├── board_dto.go
│ │ └── git_dto.go
│
├── infra/ # 🧱 Deployment & DevOps layer
│ ├── docker/
│ │ ├── Dockerfile.gateway
│ │ ├── Dockerfile.trello
│ │ ├── Dockerfile.git
│ │ └── Dockerfile.userauth
│ ├── docker-compose.yml
│ ├── k8s/
│ │ ├── gateway-deployment.yaml
│ │ ├── trello-deployment.yaml
│ │ ├── kafka-deployment.yaml
│ │ └── ...
│ └── monitoring/
│ ├── prometheus.yml
│ └── grafana/
│
├── go.work # connects all go modules (monorepo)
├── Makefile # build + test + deploy commands
└── README.md

---

# Design Principle

| Concept                        | Description                                                  |
| ------------------------------ | ------------------------------------------------------------ |
| **Independent services**       | Each microservice has its own `go.mod`, DB, and routes.      |
| **API Gateway**                | Central entry point for auth, routing, and RBAC.             |
| **Shared Library (`/shared`)** | Common code reused by all services (logging, config, Kafka). |
| **Message Broker (Kafka)**     | Event-driven sync between Trello ↔ Git ↔ Notification.       |
| **Databases**                  | Separate PostgreSQL schema per service.                      |
| **Cache Layer (Redis)**        | Used by gateway and Trello service for fast lookups.         |
| **Monitoring Stack**           | Prometheus + Grafana integrated under `/infra/monitoring`.   |
| **Deployment Ready**           | Docker + Kubernetes manifests for each service.              |

---

# Tech Stacks

| Layer            | Tools                        |
| ---------------- | ---------------------------- |
| API Gateway      | Gin (Go)                     |
| Backend Services | Go (Gin / Chi / Fiber)       |
| Message Broker   | Kafka                        |
| Cache            | Redis                        |
| Search           | Elasticsearch                |
| DB               | PostgreSQL (per service)     |
| Frontend         | React + Vite + Tailwind      |
| Auth             | JWT + OAuth (GitHub, GitLab) |
| Infra            | Docker, Kubernetes, Helm     |
| CI/CD            | GitHub Actions / ArgoCD      |
| Monitoring       | Prometheus + Grafana         |
