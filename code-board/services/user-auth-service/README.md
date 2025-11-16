# auth

codeboard-backend/
│
├── cmd/
│ └── server/
│ └── main.go # entry point
│
├── config/
│ └── config.yaml # app configs (DB, ports, etc.)
│
├── internal/
│ ├── handler/ # HTTP layer
│ │ ├── board_handler.go
│ │ ├── card_handler.go
│ │ └── user_handler.go
│ │
│ ├── service/ # business logic
│ │ ├── board_service.go
│ │ ├── card_service.go
│ │ └── user_service.go
│ │
│ ├── repository/ # data layer
│ │ ├── board_repository.go
│ │ ├── card_repository.go
│ │ └── user_repository.go
│ │
│ ├── model/ # GORM models / DTOs
│ │ ├── board.go
│ │ ├── card.go
│ │ └── user.go
│ │
│ ├── router/ # all routes
│ │ └── router.go
│ │
│ ├── db/ # DB connection and migrations
│ │ ├── connect.go
│ │ └── migrate.go
│ │
│ ├── middleware/ # JWT, logging, CORS, etc.
│ │ ├── auth.go
│ │ └── cors.go
│ │
│ └── utils/ # helper functions (response, logger, etc.)
│ ├── response.go
│ └── logger.go
│
├── pkg/ # optional shared packages (like validations)
│ └── validator/
│ └── validator.go
│
├── .env
├── go.mod
├── go.sum
└── Dockerfile

## apis

⚙️ Auth Service — Complete Route Structure
🔐 1. Authentication Routes (/api/v1/auth)
Method Endpoint Description
POST /api/v1/auth/register Register a new user (email/password)
POST /api/v1/auth/login Login with email/password
POST /api/v1/auth/refresh Refresh JWT access token using refresh token
POST /api/v1/auth/logout Invalidate current session or refresh token
GET /api/v1/auth/me Get the currently logged-in user info (JWT protected)

🧠 2. OAuth Routes (/api/v1/oauth)
Method Endpoint Description
GET /api/v1/oauth/google/login Redirect user to Google login
GET /api/v1/oauth/google/callback Handle Google OAuth callback
GET /api/v1/oauth/github/login Redirect user to GitHub login
GET /api/v1/oauth/github/callback Handle GitHub OAuth callback

🔑 3. Password & Verification Routes (/api/v1/password)
Method Endpoint Description
POST /api/v1/password/forgot Send password reset email/token
POST /api/v1/password/reset Reset password using reset token
POST /api/v1/email/verify Send email verification link
GET /api/v1/email/verify/:token Verify email address

👤 4. User Management Routes (/api/v1/users)
Method Endpoint Description
GET /api/v1/users/:id Get user by ID (admin protected)
GET /api/v1/users List all users (admin protected)
PUT /api/v1/users/:id Update user details (admin/self)
DELETE /api/v1/users/:id Delete user (admin only)

🛡️ 5. Session & Token Routes (/api/v1/sessions)
Method Endpoint Description
GET /api/v1/sessions List active sessions for user
DELETE /api/v1/sessions/:id Revoke a session/token
DELETE /api/v1/sessions/all Revoke all sessions (logout all devices)
