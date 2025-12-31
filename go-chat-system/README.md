Real-Time Chat Application
A production-ready, scalable real-time chat application built with Go, featuring WebSocket communication, WebRTC calling, and clean architecture principles.
Features

✅ Real-time 1-to-1 and group messaging
✅ Online/Offline presence tracking
✅ Last seen timestamps
✅ Message delivery and read receipts (double tick)
✅ Audio & Video calling with WebRTC
✅ JWT authentication
✅ Horizontal scalability with gRPC
✅ Clean Architecture (Domain-Driven Design)
✅ PostgreSQL for persistence
✅ Redis for presence and caching

Architecture
┌─────────────────────────────────────────────┐
│           Client (Web/Mobile)               │
│  REST API | WebSocket | WebRTC              │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│          Gateway Server (Go)                │
│  HTTP Handlers | WS Hub | gRPC Server       │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│           Core Services                     │
│  Auth | Chat | Presence | Signaling         │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│             Storage Layer                   │
│  PostgreSQL (Persistent) | Redis (Cache)    │
└─────────────────────────────────────────────┘
Project Structure
chat-app/
├── cmd/server/              # Application entry point
├── internal/
│   ├── config/              # Configuration management
│   ├── domain/              # Core business entities & interfaces
│   ├── service/             # Business logic layer
│   ├── repository/          # Data access implementations
│   │   ├── postgres/        # PostgreSQL repositories
│   │   └── redis/           # Redis repositories
│   └── transport/           # Input adapters
│       ├── http/            # REST API handlers
│       ├── websocket/       # WebSocket engine
│       └── grpc/            # gRPC server
├── pkg/                     # Public libraries
│   ├── logger/              # Structured logging
│   └── pb/                  # Generated gRPC code
├── api/proto/               # Protocol buffer definitions
├── migrations/              # Database migrations
└── deploy/                  # Deployment configs
Prerequisites

Go 1.21 or higher
PostgreSQL 15+
Redis 7+
Docker & Docker Compose (optional)
Protocol Buffers compiler (for gRPC development)

Quick Start
Option 1: Using Docker Compose (Recommended)
bash# Clone the repository
git clone <repository-url>
cd chat-app

# Start all services
make docker-up

# View logs
make docker-logs

# Stop services
make docker-down
The application will be available at:

HTTP API: http://localhost:8080
WebSocket: ws://localhost:8080/api/ws
gRPC: localhost:9090

Option 2: Local Development
bash# 1. Install dependencies
go mod download

# 2. Install development tools
make install-tools

# 3. Setup environment variables
cp .env.example .env
# Edit .env with your configuration

# 4. Start PostgreSQL and Redis
docker run -d --name postgres -p 5432:5432 \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=chatapp \
  postgres:15-alpine

docker run -d --name redis -p 6379:6379 redis:7-alpine

# 5. Run migrations
make migrate-up

# 6. Generate gRPC code (if modified proto files)
make gen-proto

# 7. Run the application
make run
API Documentation
Authentication
Register
bashcurl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "securepassword123",
    "display_name": "John Doe"
  }'
Login
bashcurl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "securepassword123"
  }'
Response:
json{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": "abc123",
    "username": "john_doe",
    "email": "john@example.com",
    "display_name": "John Doe"
  }
}
Conversations
Create Conversation
bashcurl -X POST http://localhost:8080/api/conversations \
  -H "Authorization: Bearer <token>" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "one_to_one",
    "member_ids": ["user_id_2"]
  }'
Get Messages
bashcurl -X GET "http://localhost:8080/api/conversations/{id}/messages?limit=50&offset=0" \
  -H "Authorization: Bearer <token>"
WebSocket Connection
javascriptconst token = "your_jwt_token";
const ws = new WebSocket(`ws://localhost:8080/api/ws?token=${token}`);

ws.onopen = () => {
  console.log("Connected to WebSocket");
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log("Received:", message);
};

// Send a message
ws.send(JSON.stringify({
  type: "chat_message",
  payload: {
    conversation_id: "conv_123",
    content: "Hello, World!",
    message_type: "text"
  }
}));
WebSocket Message Types
Client → Server

Send Message

json{
  "type": "chat_message",
  "payload": {
    "conversation_id": "conv_123",
    "content": "Hello!",
    "message_type": "text"
  }
}

Typing Indicator

json{
  "type": "typing",
  "payload": {
    "conversation_id": "conv_123",
    "is_typing": true
  }
}

Read Receipt

json{
  "type": "read_receipt",
  "payload": {
    "message_id": "msg_123"
  }
}

Call Offer

json{
  "type": "call_offer",
  "payload": {
    "call_id": "call_123",
    "recipient_id": "user_456",
    "offer": "sdp_offer_string"
  }
}
Server → Client

New Message

json{
  "type": "new_message",
  "payload": {
    "id": "msg_123",
    "conversation_id": "conv_123",
    "sender_id": "user_456",
    "content": "Hello!",
    "message_type": "text",
    "created_at": "2024-01-01T12:00:00Z"
  }
}

Presence Update

json{
  "type": "presence_update",
  "payload": {
    "user_id": "user_456",
    "status": "online",
    "timestamp": 1234567890
  }
}
WebRTC Calling Flow

Caller sends call_offer via WebSocket
Server validates and forwards to recipient
Callee receives offer and sends call_answer
Both parties exchange ICE candidates
Media flows peer-to-peer (not through server)

Database Schema
Key tables:

users - User accounts
conversations - Chat rooms (1-to-1 or group)
conversation_members - Membership tracking
messages - Message history
message_delivery - Delivery & read status
calls - Call history

See migrations/ for complete schema.
Scalability
Single Node

In-memory WebSocket hub
Redis for presence
PostgreSQL for persistence

Multi-Node

gRPC for inter-node communication
Redis Pub/Sub for message fan-out
Sticky sessions for WebSocket connections

Production Considerations

Load balancer with WebSocket support
Database connection pooling
Redis cluster for high availability
Monitoring and observability
Rate limiting and DDoS protection

Development
Run Tests
bashmake test
Create Migration
bashmake migrate-create name=add_user_status
Build Binary
bashmake build
Clean Build Artifacts
bashmake clean
Environment Variables
VariableDescriptionDefaultHTTP_PORTHTTP server port8080GRPC_PORTgRPC server port9090DB_HOSTPostgreSQL hostlocalhostDB_PORTPostgreSQL port5432DB_USERDatabase userpostgresDB_PASSWORDDatabase passwordpostgresDB_NAMEDatabase namechatappREDIS_HOSTRedis hostlocalhostREDIS_PORTRedis port6379JWT_SECRETJWT signing secretchange-in-production
Security Considerations

Use strong JWT secrets in production
Enable TLS for WebSocket connections (wss://)
Implement rate limiting
Validate all user inputs
Use prepared statements (already implemented)
Enable CORS only for trusted origins
Store passwords with bcrypt (already implemented)

Performance Optimization

Database indexes on frequently queried columns
Connection pooling for PostgreSQL
Redis caching for hot data
WebSocket message batching
Horizontal scaling with gRPC

Monitoring
Recommended tools:

Prometheus for metrics
Grafana for dashboards
Jaeger for distributed tracing
ELK stack for log aggregation

Contributing

Fork the repository
Create a feature branch
Make your changes
Write tests
Submit a pull request

License
MIT License - see LICENSE file for details
Support
For issues and questions:

Open an issue on GitHub
Check existing documentation
Review the architecture document