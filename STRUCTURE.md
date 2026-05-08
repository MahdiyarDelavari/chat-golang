# 🏗️ Project Structure & File Guide

This document provides a detailed breakdown of the application's structure and what each file is responsible for. 

## 🗺️ High-Level Layout

* **📂 `backend/`** - Contains the entire Go-based REST API and WebSocket server.
* **📄 `README.md`** - General overview, tech stack, and run instructions.
* **📄 `STRUCTURE.md`** - This directory map (you are here).

---

## 🖥️ Backend Directory Breakdown (`/backend`)

The backend follows a standard Go project layout, keeping the entry point in `cmd` and the application logic hidden and encapsulated inside `internal`.

### 1️⃣ Root Files

- `go.mod` / `go.sum` - Defines Go module versions and dependency checksums.
- `rest.http` - A playground file containing API requests you can trigger directly via the "REST Client" VS Code extension.

### 2️⃣ Application Entry Point

- **`cmd/api/main.go`**
  The core starting point of the application. It loads configurations, initializes the SQLite database, creates the WebSocket hub, registers middleware/routes, sets up the HTTP server, and handles graceful shutdowns.

### 3️⃣ Configuration & Databases

- **`config/dev.env`**
  Environment variables for local development (database path, secure JWT keys, server ports).
- **`sqlite/dev/`**
  Locally stored SQLite database files.

### 4️⃣ Internal Domain Logic (`/internal`)

The `internal` directory limits imports to this module, ensuring clean architecture.

#### 📂 `internal/config/`

- `config.go` - Reads `dev.env` and maps text-based environment variables into structured Go objects using `cleanenv` or similar decoding.

#### 📂 `internal/db/`

- `db.go` - Opens the SQLite connection pool, pings the database, and provides a global access point for models to query against.

#### 📂 `internal/middlewares/`

- `authenticate.go` - Intercepts incoming requests on protected routes to verify the JWT Access Token.
- `authenticateHandler.go` - Extracts bearer tokens from headers and attaches user data to the request context.
- `cors.go` - Appends "Cross-Origin Resource Sharing" headers to requests, enabling the Angular frontend to consume the API across different local ports.
- `muxlogger.go` - A wrapper that logs incoming HTTP requests (method, URI, execution time) to the console.

#### 📂 `internal/models/`

Handles Data Access Object (DAO) operations. Contains SQL queries to insert/read/update the database.

- `message.go` - DB operations for chat messages.
- `private.go` - DB operations for distinct private conversations between users.
- `user.go` - Manages user accounts, fetching by email, setting up logins, and updating refresh tokens for active sessions.

#### 📂 `internal/realtime/`

Orchestrates the live WebSocket pipeline.

- `client.go` - Represents a single WebSocket connection (a user). Handles reading incoming JSON signals and writing outbound channel messages back over the TCP socket.
- `event.go` - Defines structures for different real-time events (e.g., `UserOnline`, `MessageSent`).
- `hub.go` - The global registry. It tracks all connected clients by their User ID and broadcasts events across channels to all specific users concurrently.

#### 📂 `internal/routes/`

Contains HTTP Handlers. These functions parse incoming HTTP JSON request bodies, utilize models to make database queries, and return formatted API responses.

- `auths.go` - Login, Registration, Logout, Token Renewal, and Current User fetch operations.
- `conversations.go` - Fetching lists of active chats, joining a chat room, and grabbing message histories.
- `files.go` - Endpoints for handling multipart-form file uploads and delivering binary file downloads.
- `healths.go` - Small functions to ping the server to ensure its HTTP & WebSocket routers are alive.
- `routes.go` - The master router file that wires up all endpoints from this directory into the multiplexer (ServeMux) and binds the middlewares.
- `users.go` - Endpoints to search for and fetch users.
- `websockets.go` - The endpoint (`/api/ws`) that upgrades a standard HTTP GET Request into a persistent TCP WebSocket connection, linking it securely to a User ID.

#### 📂 `internal/utils/`

Stateless helper functions.

- `apiresponse.go` - Standardizes JSON responses uniformly (success vs error shapes).
- `jwt.go` - Contains logic to encode secure Access Tokens using HMAC keys and parse incoming tokens.
- `passwordhash.go` - Wraps Bcrypt logic for hashing passwords before DB entry, and verifying plain passwords against stored hashes.
- `refreshtoken.go` - Utilities for generating secure, cryptographically random refresh tokens for long-lived login sessions (web vs mobile multiplexing).
