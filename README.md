<div align="center">
  <h1>💬 Golang Chat Application</h1>

  ![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
  ![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)
</div>

<br/>

It supports user registration, JWT-based authentication, private conversations, real-time messaging using WebSockets, and file sharing.

## 🛠️ Tech Stack

### Backend

- **Language:** Go (1.25)
- **Router:** Standard Library `net/http` ServeMux
- **Database:** SQLite (`modernc.org/sqlite`)
- **Real-time Engine:** WebSockets (`github.com/coder/websocket`)
- **Authentication:** JWT (`golang-jwt/jwt/v5`) with support for both Web and Mobile platforms (refresh tokens)

## 📂 Directory Structure (Backend)

```text
backend/
├── cmd/
│   └── api/
│       └── main.go                 # App entry point, starts the HTTP server & WebSocket hub
├── config/
│   └── dev.env                     # Environment variables (e.g., JWT_KEY, DB_PATH)
├── internal/
│   ├── config/                     # Configuration schemas and env loading logic
│   ├── db/                         # Database connection and initialization
│   ├── middlewares/                # Custom HTTP Middlewares (CORS, Logger, JWT Auth)
│   ├── models/                     # DB schemas & access logic (User, Message, Private Conversations)
│   ├── realtime/                   # WebSocket events mapping, hub registry, and client management
│   ├── routes/                     # HTTP Handlers and route definitions
│   └── utils/                      # Helper tools (API responses, Password Hashing, JWT tools)
├── sqlite/
│   └── dev/                        # SQLite development database location
├── rest.http                       # API test file for HTTP clients (VS Code REST Client)
└── go.mod                          # Go module dependencies
```

## ✨ Features

- **Authentication:** Allows users to sign up via email and password. Generates Access & Refresh JWT tokens. Refresh tokens support multiple platforms (Web/Mobile multiplexing via `X-Platform` header constraints).
- **Real-Time Communication:** WebSockets allow users to detect when others go online, send messages smoothly, and retrieve queued, undelivered messages simultaneously upon connecting.
- **Conversations:** Supports retrieving users and connecting to private chat rooms (`privates`). Messages are stored within SQL using pagination.
- **File Sharing:** Complete standard file-sharing abilities with endpoints built for uploading and downloading content scoped per conversation.
- **Monitoring:** Standard HTTP & WebSocket Health Check routes built in to observe platform health metrics.

## 🔌 API Endpoints Overview

For a more comprehensive view of API limits, refer to the included `rest.http` file which has pre-configured payload checks.

- **Health**:
  - `GET /api/health-check-http`
  - `GET ws://{url}/api/health-check-ws`
- **Auth**:
  - `POST /api/auth/register-email`
  - `POST /api/auth/login-email`
  - `POST /api/auth/logout` _(Requires Auth)_
  - `POST /api/auth/refresh-session` _(Requires Auth)_
  - `GET /api/auth/current-user` _(Requires Auth)_
- **Users**:
  - `GET /api/users/{id}` _(Requires Auth)_
- **Conversations**:
  - `GET /api/conversations` _(Requires Auth)_
  - `POST /api/conversations/privates/join` _(Requires Auth)_
  - `GET /api/conversations/privates/{private_id}` _(Requires Auth)_
  - `GET /api/conversations/privates/{private_id}/messages?page=...&limit=...` _(Requires Auth)_
- **Files**:
  - `POST /api/files/{private_id}` _(Requires Auth)_
  - `GET /api/files/` _(Requires Auth)_
- **WebSockets**:
  - `GET ws://{url}/api/ws` _(Requires Auth token)_

## 🚀 How to Run

1. Navigate to the backend folder:
   ```sh
   cd backend
   ```
2. Build and run the server:
   ```sh
   go run ./cmd/api/main.go
   ```
3. The server uses default settings typically hosted at `http://localhost:8082` (or defined in `dev.env`).
